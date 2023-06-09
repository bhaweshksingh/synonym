package middleware

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"synonym/pkg/http/internal/resperr"
	"synonym/pkg/http/internal/utils"
)

func WithErrorHandler(lgr *zap.Logger, next func(resp http.ResponseWriter, req *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := next(w, r)
		if err == nil {
			return
		}

		lgr.Error(err.Error())

		utils.WriteFailureResponse(w, resperr.NewResponseError(http.StatusInternalServerError, fmt.Sprintf("could not process the request: %q", err.Error())))
	}
}

func WithReqResLog(lgr *zap.Logger, next func(resp http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		reqBody, _ := ioutil.ReadAll(req.Body)

		bd := ioutil.NopCloser(bytes.NewBuffer(reqBody))
		req.Body = bd

		respWriter := utils.NewCopyWriter(resp)

		next(respWriter, req)
	}
}

func WithSecurityHeaders(next func(resp http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Referrer-Policy", "strict-origin")
		resp.Header().Set("Content-Security-Policy", "script-src 'self'")
		resp.Header().Set("Strict-Transport-Security", "max-age=31536000 ; includeSubDomains")
		resp.Header().Set("X-Content-Type-Options", "nosniff")

		next(resp, req)
	}
}
