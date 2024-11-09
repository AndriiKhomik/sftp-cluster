package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	hosts := strings.Split(os.Getenv("GO_HOSTS"), ",")
	port := os.Getenv("GO_HOST_PORT")
	username := os.Getenv("GO_SERVER_USERNAME")
	password := os.Getenv("GO_SERVER_PASSWORD")
	remoteDir := os.Getenv("GO_SERVER_UPLOADS_PATH")

	for _, host := range hosts {
		sftpCllient, err := createSFTPClient(host, port, username, password)
		if err != nil {
			log.Printf("Failed to connect for %s: %v", host, err)
			continue
		}

		getDataFromSftp(sftpCllient, remoteDir)

		sftpCllient.Close()
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/sftp/", detailedHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failded to start server: %v", err)
	}
}
