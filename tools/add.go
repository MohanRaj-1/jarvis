package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AddInput struct {
	A int `json:"first_number" jsonschema:"first number to add"`
	B int `json:"second_number" jsonschema:"second number to add"`
}

type AddOutput struct {
	Result int `json:"result"`
}

func Add(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in AddInput,
) (*mcp.CallToolResult, AddOutput, error) {

	return nil, AddOutput{
		Result: in.A + in.B,
	}, nil
}
