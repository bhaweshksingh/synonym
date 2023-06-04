package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"net/http"
)

func ParseRequest(req *http.Request, data interface{}) (err error) {
	if req == nil {
		return fmt.Errorf("ParseRequest. Error: request is nil")
	}

	if req.Body == nil {
		return fmt.Errorf("ParseRequest. Error: request body  is nil")
	}

	decoder := json.NewDecoder(req.Body)

	err = decoder.Decode(&data)
	if err != nil {
		return fmt.Errorf("ParseRequest.decoder.Decode. Error: %v", err)
	}

	return
}

func ParseQueryParams(req *http.Request, data interface{}) error {
	if req == nil {
		return fmt.Errorf("ParseQueryParams. Error: request is nil")
	}

	decoder := schema.NewDecoder()

	err := decoder.Decode(data, req.URL.Query())
	if err != nil {
		return fmt.Errorf("ParseQueryParams.decoder.Decode. Error: %v", err)
	}
	return nil
}
