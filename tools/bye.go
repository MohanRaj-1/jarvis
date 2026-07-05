package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ByeInput struct {
	Name string `json:"name" jsonschema:"The name of the person for farewell"`
}

type ByeOutput struct {
	Farewell string `json:"farewell"`
}

func Bye(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ByeInput,
) (*mcp.CallToolResult, ByeOutput, error) {

	return nil, ByeOutput{
		Farewell: "Bye " + in.Name + " 👋",
	}, nil
}
