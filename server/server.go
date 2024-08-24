package server

import (
	"bytes"
	"dataTransferServer/util"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

type Server struct {
}

func (s *Server) Start() {
	// Setup Server
	outFolder := "received"
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
	fmt.Printf("Server listening on %s\n", url)

	// Accept Connection
	for {
		con, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		tcpConn := con.(*net.TCPConn)
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(2 * time.Minute)
		go writeFile(con, outFolder)
	}
}

// writeFile receives and writes file to disk
func writeFile(con net.Conn, outFolder string) {
	defer con.Close()

	// 1. send Handshake to indicate success
	_, err := con.Write([]byte("HANDSHAKE_OK"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Handshake sent...")

	//2. receive file name
	fileNameBuf := make([]byte, 256)
	_, err = io.ReadFull(con, fileNameBuf)
	if err != nil {
		log.Fatal(err)
	}
	fileName := string(bytes.Trim(fileNameBuf, "\x00"))
	if fileName == "" {
		log.Fatal("file name is empty")
		return
	}
	fmt.Printf("File name received: %s\n", fileName)

	// 3. create the file
	file, err := os.Create(filepath.Join(outFolder, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 4. loop until End of File: receive and write
	fmt.Println("starting to receive data...")
	nBytes, err := io.Copy(file, con)
	if err != nil {
		if err == io.EOF {
			fmt.Println("file transfer successful ")
			return
		} else {
			log.Fatal(err)
		}
	}
	fmt.Printf("file transfer %d bytes\n", nBytes)

}
