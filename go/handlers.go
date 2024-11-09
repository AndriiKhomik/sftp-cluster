package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func indexHandler(writter http.ResponseWriter, request *http.Request) {
	summaryData := make(map[string]struct {
		Count int `json:"count"`
	})

	for serverName, data := range structuredData {
		summaryData[serverName] = struct {
			Count int `json:"count"`
		}{Count: data.Count}
	}

	writter.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writter).Encode(summaryData); err != nil {
		http.Error(writter, "Error converting to JSON", http.StatusInternalServerError)
		return
	}
}

func detailedHandler(writter http.ResponseWriter, request *http.Request) {
	serverID := strings.TrimPrefix(request.URL.Path, "/sftp/")
	if serverID == "" {
		http.Error(writter, "Server ID not specified", http.StatusBadRequest)
		return
	}

	data, exists := structuredData[serverID]
	if !exists {
		http.Error(writter, fmt.Sprintf("No data found server ID: %s", serverID), http.StatusNotFound)
		return
	}

	writter.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writter).Encode(map[string]interface{}{serverID: data}); err != nil {
		http.Error(writter, "Error converting to JSON", http.StatusInternalServerError)
		return
	}
}
