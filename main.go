package main

import (
	"bytes"
	"dataTransferServer/util"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

func main() {
	// Setup Server
	outPath := "received"
	ip4, err := util.GetLocalIP()
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%v:8080/", ip4)

	// Init Server
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("error starting server:", err)
		return
	}
	defer ln.Close()
	fmt.Printf("Server listening on %s", url)

	// Accept Connection
	for {
		con, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleFileTransfer(con, outPath)
	}
}

func handleFileTransfer(con net.Conn, outPath string) {
	defer con.Close()

	buf := new(bytes.Buffer)

	var fileSize int64
	err := binary.Read(con, binary.LittleEndian, &fileSize)
	if err != nil {
		log.Fatalf("Error reading file size: %v", err)
	}

	filenameBytes := make([]byte, 256)
	_, err = io.ReadFull(con, filenameBytes)
	if err != nil {
		log.Fatalf("Error reading file name: %v", err)
	}
	filename := string(bytes.Trim(filenameBytes, "\x00"))

	file, err := os.Create(filepath.Join(outPath, filename))
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 4*1024) // 4 KB buffer
	remainingBytes := fileSize

	for remainingBytes > 0 {
		chunkSize := int64(len(buffer))
		if remainingBytes < chunkSize {
			chunkSize = remainingBytes
		}

		n, err := io.ReadFull(con, buffer[:chunkSize])
		if err != nil {
			log.Fatalf("Error reading from connection: %v", err)
		}

		buf.Write(buffer[:n])

		_, err = file.Write(buf.Bytes())
		if err != nil {
			log.Fatalf("Error writing to file: %v", err)
		}

		buf.Reset()

		remainingBytes -= int64(n)
	}

	fmt.Println("File transfer completed:", filename)
}
