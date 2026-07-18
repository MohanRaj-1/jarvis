package greeting

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ByeInput contains the name to farewell.
type ByeInput struct {
	Name string `json:"name" jsonschema:"The name of the person for farewell"`
}

// ByeOutput contains the farewell message.
type ByeOutput struct {
	Farewell string `json:"farewell"`
}

// Bye returns a farewell for the supplied name.
func Bye(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ByeInput,
) (*mcp.CallToolResult, ByeOutput, error) {
	return nil, ByeOutput{
		Farewell: "Bye " + in.Name + " 👋",
	}, nil
}
