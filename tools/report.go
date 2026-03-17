package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gobenpark/chatwoot-mcp/chatwoot"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- Input types ---

type GetReportsSummaryInput struct {
	Since string `json:"since,omitempty"`
	Until string `json:"until,omitempty"`
}

type GetAgentSummaryInput struct {
	Since string `json:"since,omitempty"`
	Until string `json:"until,omitempty"`
}

type GetTeamSummaryInput struct {
	Since string `json:"since,omitempty"`
	Until string `json:"until,omitempty"`
}

type GetInboxSummaryInput struct {
	Since string `json:"since,omitempty"`
	Until string `json:"until,omitempty"`
}

type GetChannelSummaryInput struct {
	Since string `json:"since,omitempty"`
	Until string `json:"until,omitempty"`
}

func fmtFloat(v *float64) string {
	if v == nil {
		return "N/A"
	}
	return fmt.Sprintf("%.1fs", *v)
}

func formatSummaryEntry(sb *strings.Builder, label string, e chatwoot.SummaryReportEntry) {
	sb.WriteString(fmt.Sprintf("- [%d] %s\n", e.ID, label))
	sb.WriteString(fmt.Sprintf("    Conversations: %d, Resolved: %d\n", e.ConversationsCount, e.ResolvedConversationsCount))
	sb.WriteString(fmt.Sprintf("    Avg FRT: %s, Avg Resolution: %s\n", fmtFloat(e.AvgFirstResponseTime), fmtFloat(e.AvgResolutionTime)))
}

// RegisterReportTools registers report-related tools on the MCP server.
func RegisterReportTools(server *mcp.Server, client *chatwoot.Client) {

	// --- get_reports_summary ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_reports_summary",
		Description: "Get account-level report summary with metrics like avg first response time, avg resolution time, conversations count, and message counts. Provide since/until as dates (YYYY-MM-DD). Defaults to last 7 days.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetReportsSummaryInput) (*mcp.CallToolResult, any, error) {
		since, until := parseDateRange(input.Since, input.Until)
		summary, err := client.GetReportsSummary(ctx, since, until, "account")
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString("Account Report Summary\n")
		sb.WriteString(fmt.Sprintf("Period: %s to %s\n\n", time.Unix(since, 0).Format("2006-01-02"), time.Unix(until, 0).Format("2006-01-02")))
		sb.WriteString(fmt.Sprintf("Conversations: %d\n", summary.ConversationsCount))
		sb.WriteString(fmt.Sprintf("Resolutions: %d\n", summary.ResolutionsCount))
		sb.WriteString(fmt.Sprintf("Incoming messages: %d\n", summary.IncomingMessagesCount))
		sb.WriteString(fmt.Sprintf("Outgoing messages: %d\n", summary.OutgoingMessagesCount))
		sb.WriteString(fmt.Sprintf("Avg first response time: %s\n", fmtFloat(summary.AvgFirstResponseTime)))
		sb.WriteString(fmt.Sprintf("Avg resolution time: %s\n", fmtFloat(summary.AvgResolutionTime)))
		if summary.Previous != nil {
			sb.WriteString(fmt.Sprintf("\nPrevious period: %d conversations, %d resolutions\n",
				summary.Previous.ConversationsCount, summary.Previous.ResolutionsCount))
		}
		return textResult(sb.String()), nil, nil
	})

	// --- get_agent_summary ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_agent_summary",
		Description: "Get per-agent performance metrics including conversations count, response time, and resolution time. Provide since/until as dates (YYYY-MM-DD). Defaults to last 7 days.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetAgentSummaryInput) (*mcp.CallToolResult, any, error) {
		since, until := parseDateRange(input.Since, input.Until)
		entries, err := client.GetAgentSummary(ctx, since, until)
		if err != nil {
			return errorResult(err), nil, nil
		}
		// Build agent name map
		agentNames := map[int]string{}
		if agents, err := client.ListAgents(ctx); err == nil {
			for _, a := range agents {
				agentNames[a.ID] = fmt.Sprintf("%s <%s>", a.Name, a.Email)
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Agent Summary (%s to %s)\n\n", time.Unix(since, 0).Format("2006-01-02"), time.Unix(until, 0).Format("2006-01-02")))
		for _, e := range entries {
			label := agentNames[e.ID]
			if label == "" {
				label = fmt.Sprintf("Agent #%d", e.ID)
			}
			formatSummaryEntry(&sb, label, e)
		}
		if len(entries) == 0 {
			sb.WriteString("No agent data available.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- get_team_summary ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_team_summary",
		Description: "Get per-team performance metrics. Provide since/until as dates (YYYY-MM-DD). Defaults to last 7 days.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetTeamSummaryInput) (*mcp.CallToolResult, any, error) {
		since, until := parseDateRange(input.Since, input.Until)
		entries, err := client.GetTeamSummary(ctx, since, until)
		if err != nil {
			return errorResult(err), nil, nil
		}
		// Build team name map
		teamNames := map[int]string{}
		if teams, err := client.ListTeams(ctx); err == nil {
			for _, t := range teams {
				teamNames[t.ID] = t.Name
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Team Summary (%s to %s)\n\n", time.Unix(since, 0).Format("2006-01-02"), time.Unix(until, 0).Format("2006-01-02")))
		for _, e := range entries {
			label := teamNames[e.ID]
			if label == "" {
				label = fmt.Sprintf("Team #%d", e.ID)
			}
			formatSummaryEntry(&sb, label, e)
		}
		if len(entries) == 0 {
			sb.WriteString("No team data available.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- get_inbox_summary ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_inbox_summary",
		Description: "Get per-inbox performance metrics. Provide since/until as dates (YYYY-MM-DD). Defaults to last 7 days.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetInboxSummaryInput) (*mcp.CallToolResult, any, error) {
		since, until := parseDateRange(input.Since, input.Until)
		entries, err := client.GetInboxSummary(ctx, since, until)
		if err != nil {
			return errorResult(err), nil, nil
		}
		// Build inbox name map
		inboxNames := map[int]string{}
		if inboxes, err := client.ListInboxes(ctx); err == nil {
			for _, i := range inboxes {
				inboxNames[i.ID] = i.Name
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Inbox Summary (%s to %s)\n\n", time.Unix(since, 0).Format("2006-01-02"), time.Unix(until, 0).Format("2006-01-02")))
		for _, e := range entries {
			label := inboxNames[e.ID]
			if label == "" {
				label = fmt.Sprintf("Inbox #%d", e.ID)
			}
			formatSummaryEntry(&sb, label, e)
		}
		if len(entries) == 0 {
			sb.WriteString("No inbox data available.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- get_channel_summary ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_channel_summary",
		Description: "Get per-channel performance metrics grouped by channel type (email, web, api, etc.). Provide since/until as dates (YYYY-MM-DD). Defaults to last 7 days.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetChannelSummaryInput) (*mcp.CallToolResult, any, error) {
		since, until := parseDateRange(input.Since, input.Until)
		channels, err := client.GetChannelSummary(ctx, since, until)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Channel Summary (%s to %s)\n\n", time.Unix(since, 0).Format("2006-01-02"), time.Unix(until, 0).Format("2006-01-02")))
		for _, ch := range channels {
			sb.WriteString(fmt.Sprintf("- %s\n", ch.ChannelType))
			sb.WriteString(fmt.Sprintf("    Conversations: %d, Resolved: %d\n", ch.ConversationsCount, ch.ResolvedConversationsCount))
			sb.WriteString(fmt.Sprintf("    Avg FRT: %s, Avg Resolution: %s\n", fmtFloat(ch.AvgFirstResponseTime), fmtFloat(ch.AvgResolutionTime)))
		}
		if len(channels) == 0 {
			sb.WriteString("No channel data available.")
		}
		return textResult(sb.String()), nil, nil
	})
}

func parseDateRange(sinceStr, untilStr string) (int64, int64) {
	now := time.Now()
	until := now.Unix()
	since := now.AddDate(0, 0, -7).Unix()

	if sinceStr != "" {
		if t, err := time.Parse("2006-01-02", sinceStr); err == nil {
			since = t.Unix()
		}
	}
	if untilStr != "" {
		if t, err := time.Parse("2006-01-02", untilStr); err == nil {
			until = t.Unix()
		}
	}
	return since, until
}
