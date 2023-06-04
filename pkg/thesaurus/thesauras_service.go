package thesaurus

import (
	"sync"
	"synonym/pkg/thesaurus/model"
)

type Thesaurus struct {
	sync.Mutex
	wordTable       map[int]model.Word
	synonymTable    map[int][]model.Synonym
	nextID          int
	addSynonymCh    chan model.Word
	searchSynonymCh chan string
	resCh           chan []model.Word
}

type Service interface {
	AddSynonym(word model.Word)
	SearchSynonym(word model.Word) []model.Word
}

func NewThesaurus() Service {
	t := &Thesaurus{
		wordTable:       make(map[int]model.Word),
		synonymTable:    make(map[int][]model.Synonym),
		nextID:          1,
		addSynonymCh:    make(chan model.Word, 100),
		searchSynonymCh: make(chan string, 100),
		resCh:           make(chan []model.Word, 100),
	}
	go t.run()
	return t
}

func (t *Thesaurus) run() {
	for {
		select {
		case word := <-t.addSynonymCh:
			t.Lock()
			id := t.nextID
			t.nextID++
			t.wordTable[id] = word
			for _, synonym := range word.Synonyms {
				t.synonymTable[id] = append(t.synonymTable[id], synonym)
				t.synonymTable[synonym.ID] = append(t.synonymTable[synonym.ID], model.Word{ID: id})
			}
			t.Unlock()
		case word := <-t.searchSynonymCh:
			t.Lock()
			for _, w := range t.wordTable {
				if w.Word == word {
					t.resCh <- t.synonymTable[w.ID]
					t.Unlock()
					break
				}
			}
			t.resCh <- nil
			t.Unlock()
		}
	}
}

func (t *Thesaurus) AddSynonym(word model.Word) {
	t.addSynonymCh <- word
}

func (t *Thesaurus) SearchSynonym(word model.Word) []model.Word {
	t.searchSynonymCh <- word.Word
	synonyms := <-t.resCh
	return synonyms
}
