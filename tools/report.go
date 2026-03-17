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
		sb.WriteString(fmt.Sprintf("Avg first response time: %.1fs\n", summary.AvgFirstResponseTime))
		sb.WriteString(fmt.Sprintf("Avg resolution time: %.1fs\n", summary.AvgResolutionTime))
		return textResult(sb.String()), nil, nil
	})

	// --- get_agent_summary ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_agent_summary",
		Description: "Get per-agent performance metrics including conversations count, response time, and resolution time. Provide since/until as dates (YYYY-MM-DD). Defaults to last 7 days.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetAgentSummaryInput) (*mcp.CallToolResult, any, error) {
		since, until := parseDateRange(input.Since, input.Until)
		agents, err := client.GetAgentSummary(ctx, since, until)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Agent Summary (%s to %s)\n\n", time.Unix(since, 0).Format("2006-01-02"), time.Unix(until, 0).Format("2006-01-02")))
		for _, a := range agents {
			sb.WriteString(fmt.Sprintf("- [%d] %s <%s>\n", a.ID, a.Name, a.Email))
			sb.WriteString(fmt.Sprintf("    Conversations: %d, Resolutions: %d\n", a.ConversationsCount, a.ResolutionsCount))
			sb.WriteString(fmt.Sprintf("    Avg FRT: %.1fs, Avg Resolution: %.1fs\n", a.AvgFirstResponseTime, a.AvgResolutionTime))
			sb.WriteString(fmt.Sprintf("    Messages in: %d, out: %d\n", a.IncomingMessagesCount, a.OutgoingMessagesCount))
		}
		if len(agents) == 0 {
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
		teams, err := client.GetTeamSummary(ctx, since, until)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Team Summary (%s to %s)\n\n", time.Unix(since, 0).Format("2006-01-02"), time.Unix(until, 0).Format("2006-01-02")))
		for _, t := range teams {
			sb.WriteString(fmt.Sprintf("- [%d] %s\n", t.ID, t.Name))
			sb.WriteString(fmt.Sprintf("    Conversations: %d, Resolutions: %d\n", t.ConversationsCount, t.ResolutionsCount))
			sb.WriteString(fmt.Sprintf("    Avg FRT: %.1fs, Avg Resolution: %.1fs\n", t.AvgFirstResponseTime, t.AvgResolutionTime))
			sb.WriteString(fmt.Sprintf("    Messages in: %d, out: %d\n", t.IncomingMessagesCount, t.OutgoingMessagesCount))
		}
		if len(teams) == 0 {
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
		inboxes, err := client.GetInboxSummary(ctx, since, until)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Inbox Summary (%s to %s)\n\n", time.Unix(since, 0).Format("2006-01-02"), time.Unix(until, 0).Format("2006-01-02")))
		for _, i := range inboxes {
			sb.WriteString(fmt.Sprintf("- [%d] %s\n", i.ID, i.Name))
			sb.WriteString(fmt.Sprintf("    Conversations: %d, Resolutions: %d\n", i.ConversationsCount, i.ResolutionsCount))
			sb.WriteString(fmt.Sprintf("    Avg FRT: %.1fs, Avg Resolution: %.1fs\n", i.AvgFirstResponseTime, i.AvgResolutionTime))
			sb.WriteString(fmt.Sprintf("    Messages in: %d, out: %d\n", i.IncomingMessagesCount, i.OutgoingMessagesCount))
		}
		if len(inboxes) == 0 {
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
			sb.WriteString(fmt.Sprintf("    Conversations: %d, Resolutions: %d\n", ch.ConversationsCount, ch.ResolutionsCount))
			sb.WriteString(fmt.Sprintf("    Avg FRT: %.1fs, Avg Resolution: %.1fs\n", ch.AvgFirstResponseTime, ch.AvgResolutionTime))
			sb.WriteString(fmt.Sprintf("    Messages in: %d, out: %d\n", ch.IncomingMessagesCount, ch.OutgoingMessagesCount))
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
