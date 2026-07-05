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
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "bye",
			Description: "Farewell a person",
		},
		tools.Bye,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "add",
			Description: "Add two numbers",
		},
		tools.Add,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "weeday",
			Description: "Gives the current day of the week",
		},
		tools.WeekDay,
	)

	log.Println("Jarvis starting...")

	if err := server.Run(
		context.Background(),
		&mcp.StdioTransport{},
	); err != nil {
		log.Fatal(err)
	}
}
