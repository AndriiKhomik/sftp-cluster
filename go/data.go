package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/sftp"
)

var structuredData = make(map[string]struct {
	Count int      `json:"count"`
	Data  []string `json:"data"`
})

func processFile(sftpClient *sftp.Client, remoteFilePath string) {
	file, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		log.Printf("Failed to open file %s: %v", remoteFilePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Text()
		parts := strings.Fields(data)

		if len(parts) < 3 {
			log.Println("Unexpected data format:", data)
			continue
		}

		timestamp := fmt.Sprintf("%s %s", parts[0], strings.TrimSuffix(parts[1], ","))
		user := parts[2]

		if _, exists := structuredData[user]; !exists {
			structuredData[user] = struct {
				Count int      `json:"count"`
				Data  []string `json:"data"`
			}{Count: 0, Data: []string{}}
		}

		entry := structuredData[user]
		entry.Count++
		entry.Data = append(entry.Data, timestamp)
		structuredData[user] = entry
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file %s: %v", remoteFilePath, err)
	}
}
