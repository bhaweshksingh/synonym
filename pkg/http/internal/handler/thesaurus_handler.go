package handler

import (
	"fmt"
	"net/http"
	"synonym/pkg/http/contract"
	"synonym/pkg/http/internal/utils"
	"synonym/pkg/thesaurus"
	"synonym/pkg/thesaurus/dto"

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
	var synonymAddReq dto.AddSynonymsReq
	err := utils.ParseRequest(req, &synonymAddReq)
	if err != nil {
		return err
	}

	t.ThesaurusSvc.AddSynonym(synonymAddReq.Words)

	t.lgr.Debug("msg", zap.String("eventCode", utils.WordAdded))
	utils.WriteSuccessResponse(resp, http.StatusCreated, contract.WordsAdditionSuccess)
	return nil
}

func (t *ThesaurusHandler) SearchSynonyms(resp http.ResponseWriter, req *http.Request) error {

	var wordSearchReq dto.WordSearchReq

	err := utils.ParseQueryParams(req, &wordSearchReq)
	if err != nil {
		return err
	}

	if len(wordSearchReq.Word) == 0 {
		return fmt.Errorf("empty word")
	}

	synonyms := t.ThesaurusSvc.SearchSynonym(wordSearchReq.Word)

	if synonyms == nil {
		return fmt.Errorf("Synonyms not found")
	}

	utils.WriteSuccessResponse(resp, http.StatusOK, dto.SynonymsResponse{Synonyms: synonyms})
	return nil
}
