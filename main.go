package main

import (
	"dataTransferServer/server"
)

func main() {
	fileServer := server.Server{}
	fileServer.Start()
}
