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

type ListCannedResponsesInput struct{}

type SearchCannedResponsesInput struct {
	Search string `json:"search"`
}

type CreateCannedResponseInput struct {
	ShortCode string `json:"short_code"`
	Content   string `json:"content"`
}

type DeleteCannedResponseInput struct {
	ID int `json:"id"`
}

type ListCustomAttributesInput struct{}

type ListCustomFiltersInput struct {
	FilterType string `json:"filter_type,omitempty"`
}

type ListAutomationRulesInput struct{}

type ListWebhooksInput struct{}

type CreateWebhookInput struct {
	URL           string   `json:"url"`
	Subscriptions []string `json:"subscriptions,omitempty"`
}

type DeleteWebhookInput struct {
	WebhookID int `json:"webhook_id"`
}

// RegisterAutomationTools registers canned responses, custom attributes, custom filters,
// automation rules, and webhook tools.
func RegisterAutomationTools(server *mcp.Server, client *chatwoot.Client) {

	// --- list_canned_responses ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_canned_responses",
		Description: "List all canned responses (saved reply templates) in the account.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListCannedResponsesInput) (*mcp.CallToolResult, any, error) {
		responses, err := client.ListCannedResponses(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, r := range responses {
			sb.WriteString(fmt.Sprintf("- [%d] /%s — %s\n", r.ID, r.ShortCode, truncate(r.Content, 80)))
		}
		if sb.Len() == 0 {
			sb.WriteString("No canned responses found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- create_canned_response ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_canned_response",
		Description: "Create a new canned response (saved reply template). short_code is the shortcut, content is the template text.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateCannedResponseInput) (*mcp.CallToolResult, any, error) {
		cr, err := client.CreateCannedResponse(ctx, chatwoot.CreateCannedResponseRequest{
			ShortCode: input.ShortCode,
			Content:   input.Content,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Canned response created! ID: %d, Short code: /%s", cr.ID, cr.ShortCode)), nil, nil
	})

	// --- delete_canned_response ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_canned_response",
		Description: "Delete a canned response by ID.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteCannedResponseInput) (*mcp.CallToolResult, any, error) {
		if err := client.DeleteCannedResponse(ctx, input.ID); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Canned response #%d deleted.", input.ID)), nil, nil
	})

	// --- list_custom_attributes ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_custom_attributes",
		Description: "List all custom attribute definitions (for conversations and contacts).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListCustomAttributesInput) (*mcp.CallToolResult, any, error) {
		attrs, err := client.ListCustomAttributes(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, a := range attrs {
			model := "conversation"
			if a.AttributeModel == 1 {
				model = "contact"
			}
			sb.WriteString(fmt.Sprintf("- [%d] %s (%s, %s) — key: %s — %s\n",
				a.ID, a.AttributeDisplayName, a.AttributeDisplayType, model, a.AttributeKey, a.AttributeDescription))
		}
		if sb.Len() == 0 {
			sb.WriteString("No custom attributes defined.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- list_custom_filters ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_custom_filters",
		Description: "List saved custom filters. filter_type: 'conversation' or 'contact'.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListCustomFiltersInput) (*mcp.CallToolResult, any, error) {
		filters, err := client.ListCustomFilters(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, f := range filters {
			sb.WriteString(fmt.Sprintf("- [%d] %s (%s) — query: %s\n", f.ID, f.Name, f.FilterType, string(f.Query)))
		}
		if sb.Len() == 0 {
			sb.WriteString("No custom filters found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- list_automation_rules ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_automation_rules",
		Description: "List all automation rules configured in the account.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListAutomationRulesInput) (*mcp.CallToolResult, any, error) {
		rules, err := client.ListAutomationRules(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, r := range rules {
			active := "active"
			if !r.Active {
				active = "inactive"
			}
			conditions, _ := json.Marshal(r.Conditions)
			actions, _ := json.Marshal(r.Actions)
			sb.WriteString(fmt.Sprintf("- [%d] %s (%s, %s)\n", r.ID, r.Name, r.EventName, active))
			if r.Description != "" {
				sb.WriteString(fmt.Sprintf("    Description: %s\n", r.Description))
			}
			sb.WriteString(fmt.Sprintf("    Conditions: %s\n", string(conditions)))
			sb.WriteString(fmt.Sprintf("    Actions: %s\n", string(actions)))
		}
		if sb.Len() == 0 {
			sb.WriteString("No automation rules found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- list_webhooks ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_webhooks",
		Description: "List all webhooks configured in the account.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListWebhooksInput) (*mcp.CallToolResult, any, error) {
		webhooks, err := client.ListWebhooks(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, w := range webhooks {
			sb.WriteString(fmt.Sprintf("- [%d] %s — events: %s\n", w.ID, w.URL, strings.Join(w.Subscriptions, ", ")))
		}
		if sb.Len() == 0 {
			sb.WriteString("No webhooks configured.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- create_webhook ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_webhook",
		Description: "Create a new webhook. Subscriptions: conversation_created, conversation_status_changed, conversation_updated, message_created, message_updated, webwidget_triggered, contact_created, contact_updated.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateWebhookInput) (*mcp.CallToolResult, any, error) {
		webhook, err := client.CreateWebhook(ctx, chatwoot.CreateWebhookRequest{
			URL:           input.URL,
			Subscriptions: input.Subscriptions,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Webhook created! ID: %d, URL: %s", webhook.ID, webhook.URL)), nil, nil
	})

	// --- delete_webhook ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_webhook",
		Description: "Delete a webhook by ID.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteWebhookInput) (*mcp.CallToolResult, any, error) {
		if err := client.DeleteWebhook(ctx, input.WebhookID); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Webhook #%d deleted.", input.WebhookID)), nil, nil
	})
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
