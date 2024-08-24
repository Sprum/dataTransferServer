package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

type Client struct {
}

func (c *Client) SendFile() {
	// params
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run client.go <file_path> <ip_address>")
		return
	}
	filePath := os.Args[1]
	ipAddress := os.Args[2]
	url := fmt.Sprintf("%s:8080", ipAddress)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	// Extract filename
	filename := filepath.Base(filePath)

	// establish connection
	con, err := net.Dial("tcp", url)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer con.Close()
	fmt.Println("connection established, waiting for handshake.")

	// 1. Get Handshake
	handshakeBuf := make([]byte, len("HANDSHAKE_OK"))
	_, err = con.Read(handshakeBuf)
	if err != nil {
		log.Fatal(err)
	}
	if string(handshakeBuf) != "HANDSHAKE_OK" {
		log.Fatal("Error receiving Handshake")
	}
	fmt.Println("Handshake successful, sending filename")

	// 2. send filename
	filenameBuf := make([]byte, 256) // padding of 256 is expected by the server
	copy(filenameBuf, filename)
	_, err = con.Write(filenameBuf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Filename sent. Sending file data...")

	// 3. send file data
	_, err = io.Copy(con, file)
	if err != nil {
		log.Fatalf("Error sending file data: %v", err)
	}
	fmt.Println("File sent successfully.")
	return
}
