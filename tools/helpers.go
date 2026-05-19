package tools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}

func errorResult(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Error: %s", err.Error())},
		},
		IsError: true,
	}
}

// parseSnoozedUntil parses a string into a Unix timestamp (seconds). Accepts either
// an RFC3339 timestamp (e.g. "2026-05-23T09:00:00Z") or a Unix epoch seconds string
// (e.g. "1748000000"). Numeric input is tried first so "1748000000" is unambiguous.
func parseSnoozedUntil(s string) (int64, error) {
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		return ts, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.Unix(), nil
	}
	return 0, fmt.Errorf("must be an RFC3339 timestamp (e.g. \"2026-05-23T09:00:00Z\") or Unix epoch seconds (e.g. \"1748000000\"), got %q", s)
}
