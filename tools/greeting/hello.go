package greeting

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HelloInput contains the name to greet.
type HelloInput struct {
	Name string `json:"name" jsonschema:"The name of the person to greet"`
}

// HelloOutput contains the greeting message.
type HelloOutput struct {
	Greeting string `json:"greeting"`
}

// Hello returns a greeting for the supplied name.
func Hello(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in HelloInput,
) (*mcp.CallToolResult, HelloOutput, error) {
	return nil, HelloOutput{
		Greeting: "Hello " + in.Name + " 👋",
	}, nil
}
