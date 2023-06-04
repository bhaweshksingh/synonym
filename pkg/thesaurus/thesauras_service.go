package thesaurus

import (
	"fmt"
	"sync"
)

type Thesaurus struct {
	sync.Mutex
	adjacencyList   map[string]map[string]bool
	graphID         map[string]int
	nextGraphID     int
	addSynonymCh    chan []string
	searchSynonymCh chan string
	resCh           chan []string
}

type Service interface {
	AddSynonym(words []string)
	SearchSynonym(searchWord string) []string
}

func NewThesaurus() Service {
	t := &Thesaurus{
		adjacencyList:   make(map[string]map[string]bool),
		graphID:         make(map[string]int),
		nextGraphID:     1,
		addSynonymCh:    make(chan []string, 100),
		searchSynonymCh: make(chan string, 100),
		resCh:           make(chan []string, 100),
	}
	go t.run()
	return t
}

func (t *Thesaurus) run() {
	for {
		select {
		case words := <-t.addSynonymCh:
			t.Lock()
			fmt.Printf("\n\n Adding Words ... %q", words)
			if len(words) == 0 {
				t.Unlock()
				continue
			}

			var existingGraphID int
			for _, word := range words {
				if id, ok := t.graphID[word]; ok {
					existingGraphID = id
					break
				}
			}

			if existingGraphID == 0 {
				existingGraphID = t.nextGraphID
				t.nextGraphID++
			}

			for _, word := range words {
				if id, ok := t.graphID[word]; ok && id != existingGraphID {
					for w, gid := range t.graphID {
						if gid == id {
							t.graphID[w] = existingGraphID
						}
					}
				}

				t.graphID[word] = existingGraphID

				if _, ok := t.adjacencyList[word]; !ok {
					t.adjacencyList[word] = make(map[string]bool)
				}

				for _, otherWord := range words {
					if otherWord != word {
						t.adjacencyList[word][otherWord] = true
					}
				}
			}
			t.Unlock()
		case word := <-t.searchSynonymCh:
			t.Lock()
			graphID, ok := t.graphID[word]
			if !ok {
				t.resCh <- nil
				t.Unlock()
				continue
			}
			synonyms := []string{}

			for word, id := range t.graphID {
				if id == graphID {
					synonyms = append(synonyms, word)
				}
			}

			t.resCh <- synonyms

			t.Unlock()
		}
	}
}

func (t *Thesaurus) AddSynonym(words []string) {
	t.addSynonymCh <- words
}

func (t *Thesaurus) SearchSynonym(searchWord string) []string {
	t.searchSynonymCh <- searchWord
	synonyms := <-t.resCh
	return synonyms
}
