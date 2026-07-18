package math

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AddInput contains the two values to add.
type AddInput struct {
	A int `json:"first_number" jsonschema:"first number to add"`
	B int `json:"second_number" jsonschema:"second number to add"`
}

// AddOutput contains the sum of the input values.
type AddOutput struct {
	Result int `json:"result"`
}

// Add returns the sum of two integers.
func Add(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in AddInput,
) (*mcp.CallToolResult, AddOutput, error) {
	return nil, AddOutput{
		Result: in.A + in.B,
	}, nil
}
