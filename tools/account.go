package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/gobenpark/chatwoot-mcp/chatwoot"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- Input types ---

type ListInboxesInput struct{}
type ListAgentsInput struct{}
type ListLabelsInput struct{}
type ListTeamsInput struct{}
type GetProfileInput struct{}

type GetTeamInput struct {
	TeamID int `json:"team_id"`
}

type ListTeamMembersInput struct {
	TeamID int `json:"team_id"`
}

type ListInboxMembersInput struct {
	InboxID int `json:"inbox_id"`
}

// RegisterAccountTools registers inbox, agent, team, label, and profile tools.
func RegisterAccountTools(server *mcp.Server, client *chatwoot.Client) {

	// --- list_inboxes ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_inboxes",
		Description: "List all inboxes (channels) in the Chatwoot account.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListInboxesInput) (*mcp.CallToolResult, any, error) {
		inboxes, err := client.ListInboxes(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, inbox := range inboxes {
			status := "enabled"
			if !inbox.Enabled {
				status = "disabled"
			}
			sb.WriteString(fmt.Sprintf("- [%d] %s (%s) — %s\n", inbox.ID, inbox.Name, inbox.ChannelType, status))
		}
		if sb.Len() == 0 {
			sb.WriteString("No inboxes found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- list_agents ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_agents",
		Description: "List all agents in the Chatwoot account.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListAgentsInput) (*mcp.CallToolResult, any, error) {
		agents, err := client.ListAgents(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, agent := range agents {
			sb.WriteString(fmt.Sprintf("- [%d] %s <%s> (%s, %s)\n", agent.ID, agent.Name, agent.Email, agent.Role, agent.AvailabilityStatus))
		}
		if sb.Len() == 0 {
			sb.WriteString("No agents found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- list_labels ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_labels",
		Description: "List all labels available in the Chatwoot account.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListLabelsInput) (*mcp.CallToolResult, any, error) {
		labels, err := client.ListLabels(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, label := range labels {
			desc := label.Description
			if desc == "" {
				desc = "(no description)"
			}
			sb.WriteString(fmt.Sprintf("- [%d] %s — %s\n", label.ID, label.Title, desc))
		}
		if sb.Len() == 0 {
			sb.WriteString("No labels found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- list_teams ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_teams",
		Description: "List all teams in the Chatwoot account.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListTeamsInput) (*mcp.CallToolResult, any, error) {
		teams, err := client.ListTeams(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, team := range teams {
			desc := team.Description
			if desc == "" {
				desc = "(no description)"
			}
			sb.WriteString(fmt.Sprintf("- [%d] %s — %s\n", team.ID, team.Name, desc))
		}
		if sb.Len() == 0 {
			sb.WriteString("No teams found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- get_team ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_team",
		Description: "Get detailed information about a specific team.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetTeamInput) (*mcp.CallToolResult, any, error) {
		team, err := client.GetTeam(ctx, input.TeamID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Team #%d: %s\n", team.ID, team.Name))
		if team.Description != "" {
			sb.WriteString(fmt.Sprintf("Description: %s\n", team.Description))
		}
		return textResult(sb.String()), nil, nil
	})

	// --- list_team_members ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_team_members",
		Description: "List all agents in a specific team.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListTeamMembersInput) (*mcp.CallToolResult, any, error) {
		members, err := client.ListTeamMembers(ctx, input.TeamID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Team #%d members (%d):\n", input.TeamID, len(members)))
		for _, m := range members {
			sb.WriteString(fmt.Sprintf("- [%d] %s <%s> (%s)\n", m.ID, m.Name, m.Email, m.AvailabilityStatus))
		}
		if len(members) == 0 {
			sb.WriteString("No members in this team.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- list_inbox_members ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_inbox_members",
		Description: "List all agents assigned to a specific inbox.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListInboxMembersInput) (*mcp.CallToolResult, any, error) {
		members, err := client.ListInboxMembers(ctx, input.InboxID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Inbox #%d members (%d):\n", input.InboxID, len(members)))
		for _, m := range members {
			sb.WriteString(fmt.Sprintf("- [%d] %s <%s> (%s)\n", m.ID, m.Name, m.Email, m.AvailabilityStatus))
		}
		if len(members) == 0 {
			sb.WriteString("No agents assigned to this inbox.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- get_profile ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_profile",
		Description: "Get the authenticated user's profile information.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetProfileInput) (*mcp.CallToolResult, any, error) {
		profile, err := client.GetProfile(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Profile: %s\n", profile.Name))
		sb.WriteString(fmt.Sprintf("Email: %s\n", profile.Email))
		sb.WriteString(fmt.Sprintf("Role: %s\n", profile.Role))
		sb.WriteString(fmt.Sprintf("Account ID: %d\n", profile.AccountID))
		return textResult(sb.String()), nil, nil
	})
}
