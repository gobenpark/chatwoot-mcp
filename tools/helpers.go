package tools

import (
	"fmt"

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
