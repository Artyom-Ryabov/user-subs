package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type JSONData interface { any | []any }

type DataResponse struct {
	Data JSONData `json:"data"`
}

func SendData(w http.ResponseWriter, data JSONData, status int) error {
	res := DataResponse{Data: data}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	log.Printf("Send response - %v\n", res)
	return nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendError(w http.ResponseWriter, message string, status int, e error) {
	res := ErrorResponse{Error: message}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		msg := "Error: something went wrong with encoding json"
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(msg))
		log.Printf("%v - %v\n", msg, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	log.Printf("%v - %v\n", message, e)
}
