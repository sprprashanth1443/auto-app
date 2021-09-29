package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Info    interface{} `json:"info,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	ID      int64       `json:"id,omitempty"`
	Path    string      `json:"path,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func WriteErrorToResponse(w http.ResponseWriter, code int, _error interface{}) {
	res := Response{
		Success: false,
		Error:   _error,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Panicln(err)
	}
}

func WriteResultToResponse(w http.ResponseWriter, code int, result interface{}) {
	res := Response{
		Success: true,
		Result:  result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Panicln(err)
	}
}
