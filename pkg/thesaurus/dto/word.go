package dto

type AddSynonymsReq struct {
	Words []string `json:"synonyms" schema:"synonyms"`
}

type WordSearchReq struct {
	Word string `json:"word" schema:"word"`
}
