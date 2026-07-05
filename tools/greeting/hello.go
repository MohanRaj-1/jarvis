package greeting

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type HelloInput struct {
	Name string `json:"name" jsonschema:"The name of the person to greet"`
}

type HelloOutput struct {
	Greeting string `json:"greeting"`
}

func Hello(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in HelloInput,
) (*mcp.CallToolResult, HelloOutput, error) {
	return nil, HelloOutput{
		Greeting: "Hello " + in.Name + " 👋",
	}, nil
}
