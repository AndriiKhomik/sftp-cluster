package main

import (
	"encoding/json"
	"net/http"
)

func indexHandler(writter http.ResponseWriter, request *http.Request) {
	jsonData, err := json.Marshal(structuredData)
	if err != nil {
		http.Error(writter, "Error converting to JSON", http.StatusInternalServerError)
		return
	}
	writter.Header().Set("Content-Type", "application/json")
	writter.Write(jsonData)
}
