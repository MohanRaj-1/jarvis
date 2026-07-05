package main

import (
	"context"
	"log"

	"jarvis/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "Jarvis",
			Version: "0.1.0",
		},
		nil,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "hello",
			Description: "Greets a person",
		},
		tools.Hello,
	)

	log.Println("Jarvis starting...")

	if err := server.Run(
		context.Background(),
		&mcp.StdioTransport{},
	); err != nil {
		log.Fatal(err)
	}
}
