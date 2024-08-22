package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error on Server start:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server läuft auf Port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error on accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// get file name
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading file name:", err)
		return
	}
	filename := string(buf[:n])

	// receive and save
	outFile, err := os.Create(filepath.Join("received", filename))
	if err != nil {
		fmt.Println("Error on creating the file:", err)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, conn)
	if err != nil {
		fmt.Println("Error on writing the file:", err)
	}

	fmt.Println("File received:", filename)
}
