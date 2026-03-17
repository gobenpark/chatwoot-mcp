package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gobenpark/chatwoot-mcp/chatwoot"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- Input types ---

type ListIntegrationAppsInput struct{}

type CreateIntegrationHookInput struct {
	AppID    string `json:"app_id"`
	InboxID  *int   `json:"inbox_id,omitempty"`
	Settings string `json:"settings,omitempty"`
}

type UpdateIntegrationHookInput struct {
	HookID   int    `json:"hook_id"`
	Settings string `json:"settings,omitempty"`
}

type DeleteIntegrationHookInput struct {
	HookID int `json:"hook_id"`
}

// RegisterIntegrationTools registers integration and hook tools on the MCP server.
func RegisterIntegrationTools(server *mcp.Server, client *chatwoot.Client) {

	// --- list_integration_apps ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_integration_apps",
		Description: "List all integrations available for the account, including their hooks.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListIntegrationAppsInput) (*mcp.CallToolResult, any, error) {
		apps, err := client.ListIntegrationApps(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, app := range apps {
			enabled := "disabled"
			if app.Enabled {
				enabled = "enabled"
			}
			sb.WriteString(fmt.Sprintf("- [%s] %s (%s) — %s\n", app.ID, app.Name, app.HookType, enabled))
			if app.Description != "" {
				sb.WriteString(fmt.Sprintf("    %s\n", app.Description))
			}
			for _, hook := range app.Hooks {
				sb.WriteString(fmt.Sprintf("    Hook #%d: status=%s\n", hook.ID, hook.Status))
			}
		}
		if sb.Len() == 0 {
			sb.WriteString("No integrations available.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- create_integration_hook ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_integration_hook",
		Description: "Create a new integration hook. Requires app_id. Optional: inbox_id, settings (JSON string).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateIntegrationHookInput) (*mcp.CallToolResult, any, error) {
		createReq := chatwoot.CreateIntegrationHookRequest{
			AppID:   input.AppID,
			InboxID: input.InboxID,
		}
		if input.Settings != "" {
			createReq.Settings = json.RawMessage(input.Settings)
		}
		hook, err := client.CreateIntegrationHook(ctx, createReq)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Integration hook created! ID: %d, App: %s", hook.ID, hook.AppID)), nil, nil
	})

	// --- update_integration_hook ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_integration_hook",
		Description: "Update an integration hook's settings. Provide settings as a JSON string.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input UpdateIntegrationHookInput) (*mcp.CallToolResult, any, error) {
		updateReq := chatwoot.UpdateIntegrationHookRequest{}
		if input.Settings != "" {
			updateReq.Settings = json.RawMessage(input.Settings)
		}
		hook, err := client.UpdateIntegrationHook(ctx, input.HookID, updateReq)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Integration hook #%d updated.", hook.ID)), nil, nil
	})

	// --- delete_integration_hook ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_integration_hook",
		Description: "Delete an integration hook by ID.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteIntegrationHookInput) (*mcp.CallToolResult, any, error) {
		if err := client.DeleteIntegrationHook(ctx, input.HookID); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Integration hook #%d deleted.", input.HookID)), nil, nil
	})
}
