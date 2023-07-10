package main

import (
	"fmt"
	"github.com/eskpil/imq"
	"github.com/eskpil/imq/pkg/common"
	"github.com/eskpil/imq/pkg/protocol"
	"log"
)

func HandleTest(ctx common.Context, message protocol.Message) error {
	fmt.Println("new message on queue test")
	return nil
}

func main() {
	server, err := imq.NewServer()
	if err != nil {
		log.Fatal("could not create server", err)
	}

	if err := server.Listen("test", HandleTest); err != nil {
		log.Fatal("could not listen on queue", err)
	}

	if err := server.Start(); err != nil {
		fmt.Printf("could not start server: %v", err)
		return
	}
}
