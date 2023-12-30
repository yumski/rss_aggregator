package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	if statusCode > 499 {
		log.Printf("Responding with 5xx error: %s", message)
	}

	type errVal struct {
		Error string `json: "error"`
	}

	respondWithJson(w, statusCode, errVal{
		Error: message,
	})
}

func respondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("content-type", "application/json")

	dat, err := json.Marshal(payload)

	if err != nil {
		log.Println("Unable to marshal payload")
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(dat)
}
