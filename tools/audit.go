package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/gobenpark/chatwoot-mcp/chatwoot"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- Input types ---

type ListAuditLogsInput struct {
	Page int `json:"page,omitempty"`
}

// RegisterAuditTools registers audit log tools on the MCP server.
func RegisterAuditTools(server *mcp.Server, client *chatwoot.Client) {

	// --- list_audit_logs ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_audit_logs",
		Description: "List audit log entries showing who did what and when. Supports pagination.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListAuditLogsInput) (*mcp.CallToolResult, any, error) {
		resp, err := client.ListAuditLogs(ctx, input.Page)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Audit Logs (page %d)\n\n", input.Page))
		for _, log := range resp.Payload {
			username := "(system)"
			if log.Username != nil {
				username = *log.Username
			}
			sb.WriteString(fmt.Sprintf("- [%d] %s %s #%d by %s", log.ID, log.Action, log.AuditableType, log.AuditableID, username))
			if log.CreatedAt.Valid {
				sb.WriteString(fmt.Sprintf(" at %s", log.CreatedAt.Format("2006-01-02 15:04")))
			}
			sb.WriteString("\n")
			if len(log.AuditedChanges) > 0 && string(log.AuditedChanges) != "null" {
				sb.WriteString(fmt.Sprintf("    Changes: %s\n", string(log.AuditedChanges)))
			}
		}
		if len(resp.Payload) == 0 {
			sb.WriteString("No audit log entries found.")
		}
		return textResult(sb.String()), nil, nil
	})
}
