package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"synonym/pkg/http/internal/handler"
	"synonym/pkg/http/internal/middleware"
	"synonym/pkg/thesaurus"
)

const (
	addSynonymsPath    = "/synonyms"
	searchSynonymsPath = "/synonyms/search"
)

func NewRouter(lgr *zap.Logger, thesaurusSvc thesaurus.Service) http.Handler {
	router := mux.NewRouter()
	router.Use(handlers.RecoveryHandler())

	th := handler.NewThesaurusHandler(lgr, thesaurusSvc)

	router.HandleFunc(addSynonymsPath, withMiddlewares(lgr, middleware.WithErrorHandler(lgr, th.AddSynonyms))).Methods(http.MethodPost)
	router.HandleFunc(searchSynonymsPath, withMiddlewares(lgr, middleware.WithErrorHandler(lgr, th.SearchSynonyms))).Methods(http.MethodPost)

	return router
}

func withMiddlewares(lgr *zap.Logger, hnd http.HandlerFunc) http.HandlerFunc {
	return middleware.WithSecurityHeaders(middleware.WithReqResLog(lgr, hnd))
}
