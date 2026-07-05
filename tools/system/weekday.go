package system

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type WeekdayInput struct{}

type WeekdayOutput struct {
	Day string `json:"day"`
}

func WeekDay(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in WeekdayInput,
) (*mcp.CallToolResult, WeekdayOutput, error) {
	return nil, WeekdayOutput{
		Day: time.Now().Weekday().String(),
	}, nil
}
