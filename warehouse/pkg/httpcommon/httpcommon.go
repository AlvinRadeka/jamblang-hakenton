package httpcommon

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonErrorResponse struct {
	Error jsonError `json:"error"`
}

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"string"`
}

func ResponseJSONError(w http.ResponseWriter, code int, message string) {
	err := jsonErrorResponse{
		Error: jsonError{
			Code:    code,
			Message: message,
		},
	}

	ResponseJSON(w, code, err)
}

func ResponseJSON(w http.ResponseWriter, code int, jsonData interface{}) {
	jsonBit, err := json.Marshal(&jsonData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Response Error")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, string(jsonBit))
}
