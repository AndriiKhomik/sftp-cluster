package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func createSFTPClient(host, port, username, password string) (*sftp.Client, error) {
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect SFTP: %v", err)
	}

	return client, nil
}

func getDataFromSftp(sftpClient *sftp.Client, remoteDir string) {
	files, err := sftpClient.ReadDir(remoteDir)
	if err != nil {
		log.Printf("Failed to list directory %s: %v", remoteDir, err)
		return
	}

	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", remoteDir, file.Name())
		processFile(sftpClient, filePath)
	}
}
