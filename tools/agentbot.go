package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/gobenpark/chatwoot-mcp/chatwoot"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- Input types ---

type ListAgentBotsInput struct{}

type GetAgentBotInput struct {
	BotID int `json:"bot_id"`
}

type CreateAgentBotInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	OutgoingURL string `json:"outgoing_url"`
	BotType     string `json:"bot_type,omitempty"`
}

type UpdateAgentBotInput struct {
	BotID       int    `json:"bot_id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	OutgoingURL string `json:"outgoing_url,omitempty"`
	BotType     string `json:"bot_type,omitempty"`
}

type DeleteAgentBotInput struct {
	BotID int `json:"bot_id"`
}

// RegisterAgentBotTools registers agent bot tools on the MCP server.
func RegisterAgentBotTools(server *mcp.Server, client *chatwoot.Client) {

	// --- list_agent_bots ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_agent_bots",
		Description: "List all agent bots in the account.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListAgentBotsInput) (*mcp.CallToolResult, any, error) {
		bots, err := client.ListAgentBots(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, bot := range bots {
			sb.WriteString(fmt.Sprintf("- [%d] %s (%s)\n", bot.ID, bot.Name, bot.BotType))
			if bot.Description != "" {
				sb.WriteString(fmt.Sprintf("    Description: %s\n", bot.Description))
			}
			if bot.OutgoingURL != "" {
				sb.WriteString(fmt.Sprintf("    Outgoing URL: %s\n", bot.OutgoingURL))
			}
		}
		if sb.Len() == 0 {
			sb.WriteString("No agent bots found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- get_agent_bot ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_agent_bot",
		Description: "Get detailed information about a specific agent bot.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetAgentBotInput) (*mcp.CallToolResult, any, error) {
		bot, err := client.GetAgentBot(ctx, input.BotID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Agent Bot #%d: %s\n", bot.ID, bot.Name))
		sb.WriteString(fmt.Sprintf("Type: %s\n", bot.BotType))
		if bot.Description != "" {
			sb.WriteString(fmt.Sprintf("Description: %s\n", bot.Description))
		}
		if bot.OutgoingURL != "" {
			sb.WriteString(fmt.Sprintf("Outgoing URL: %s\n", bot.OutgoingURL))
		}
		if bot.CreatedAt.Valid {
			sb.WriteString(fmt.Sprintf("Created: %s\n", bot.CreatedAt.Format("2006-01-02 15:04")))
		}
		return textResult(sb.String()), nil, nil
	})

	// --- create_agent_bot ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_agent_bot",
		Description: "Create a new agent bot. Requires name and outgoing_url. Optional: description, bot_type.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateAgentBotInput) (*mcp.CallToolResult, any, error) {
		bot, err := client.CreateAgentBot(ctx, chatwoot.CreateAgentBotRequest{
			Name:        input.Name,
			Description: input.Description,
			OutgoingURL: input.OutgoingURL,
			BotType:     input.BotType,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Agent bot created! ID: %d, Name: %s", bot.ID, bot.Name)), nil, nil
	})

	// --- update_agent_bot ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_agent_bot",
		Description: "Update an agent bot. Provide only fields you want to change.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input UpdateAgentBotInput) (*mcp.CallToolResult, any, error) {
		updateReq := chatwoot.UpdateAgentBotRequest{}
		if input.Name != "" {
			updateReq.Name = &input.Name
		}
		if input.Description != "" {
			updateReq.Description = &input.Description
		}
		if input.OutgoingURL != "" {
			updateReq.OutgoingURL = &input.OutgoingURL
		}
		if input.BotType != "" {
			updateReq.BotType = &input.BotType
		}
		bot, err := client.UpdateAgentBot(ctx, input.BotID, updateReq)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Agent bot #%d updated! Name: %s", bot.ID, bot.Name)), nil, nil
	})

	// --- delete_agent_bot ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_agent_bot",
		Description: "Delete an agent bot by ID.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteAgentBotInput) (*mcp.CallToolResult, any, error) {
		if err := client.DeleteAgentBot(ctx, input.BotID); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Agent bot #%d deleted.", input.BotID)), nil, nil
	})
}
