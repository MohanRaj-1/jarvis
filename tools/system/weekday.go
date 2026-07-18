package system

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// WeekdayInput has no fields because the tool needs no input.
type WeekdayInput struct{}

// WeekdayOutput contains the current weekday name.
type WeekdayOutput struct {
	Day string `json:"day"`
}

// Weekday returns the current weekday name.
func Weekday(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in WeekdayInput,
) (*mcp.CallToolResult, WeekdayOutput, error) {
	return nil, WeekdayOutput{
		Day: time.Now().Weekday().String(),
	}, nil
}
