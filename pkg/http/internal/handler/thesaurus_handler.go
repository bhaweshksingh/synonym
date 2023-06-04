package handler

import (
	"fmt"
	"net/http"
	"synonym/pkg/http/contract"
	"synonym/pkg/http/internal/utils"
	"synonym/pkg/thesaurus"
	"synonym/pkg/thesaurus/model"

	"go.uber.org/zap"
)

type ThesaurusHandler struct {
	lgr          *zap.Logger
	ThesaurusSvc thesaurus.Service
}

func NewThesaurusHandler(lgr *zap.Logger, thesaurusSvc thesaurus.Service) *ThesaurusHandler {
	return &ThesaurusHandler{
		lgr:          lgr,
		ThesaurusSvc: thesaurusSvc,
	}
}

func (t *ThesaurusHandler) AddSynonyms(resp http.ResponseWriter, req *http.Request) error {
	var word model.Word
	err := utils.ParseRequest(req, &word)
	if err != nil {
		return err
	}

	t.ThesaurusSvc.AddSynonym(word)

	t.lgr.Debug("msg", zap.String("eventCode", utils.WordAdded))
	utils.WriteSuccessResponse(resp, http.StatusCreated, contract.WordAdditionSuccess)
	resp.WriteHeader(http.StatusCreated)
	return nil
}

func (t *ThesaurusHandler) SearchSynonyms(resp http.ResponseWriter, req *http.Request) error {

	var word model.Word

	err := utils.ParseQueryParams(req, &word)
	if err != nil {
		return err
	}

	synonyms := t.ThesaurusSvc.SearchSynonym(word)

	if synonyms == nil {
		return fmt.Errorf("Word not found")
	}

	utils.WriteSuccessResponse(resp, http.StatusOK, synonyms)
	return nil
}
