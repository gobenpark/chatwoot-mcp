package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gobenpark/chatwoot-mcp/chatwoot"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- Input types ---

type ListConversationsInput struct {
	Status string `json:"status,omitempty"`
	Page   int    `json:"page,omitempty"`
}

type GetConversationInput struct {
	ConversationID int `json:"conversation_id"`
}

type CreateConversationInput struct {
	InboxID    int    `json:"inbox_id"`
	ContactID  int    `json:"contact_id"`
	Message    string `json:"message,omitempty"`
	Status     string `json:"status,omitempty"`
	AssigneeID int    `json:"assignee_id,omitempty"`
	TeamID     int    `json:"team_id,omitempty"`
}

type FilterConversationsInput struct {
	Payload []FilterPayloadInput `json:"payload"`
	Page    int                  `json:"page,omitempty"`
}

type FilterPayloadInput struct {
	AttributeKey   string   `json:"attribute_key"`
	FilterOperator string   `json:"filter_operator"`
	Values         []string `json:"values"`
	QueryOperator  string   `json:"query_operator,omitempty"`
}

type GetConversationMetaInput struct{}

type UpdateConversationInput struct {
	ConversationID   int            `json:"conversation_id"`
	CustomAttributes map[string]any `json:"custom_attributes,omitempty"`
}

type GetMessagesInput struct {
	ConversationID int `json:"conversation_id"`
}

type SendMessageInput struct {
	ConversationID int    `json:"conversation_id"`
	Content        string `json:"content"`
	MessageType    string `json:"message_type,omitempty"`
	Private        bool   `json:"private,omitempty"`
}

type DeleteMessageInput struct {
	ConversationID int `json:"conversation_id"`
	MessageID      int `json:"message_id"`
}

type ToggleStatusInput struct {
	ConversationID int    `json:"conversation_id"`
	Status         string `json:"status"`
}

type TogglePriorityInput struct {
	ConversationID int    `json:"conversation_id"`
	Priority       string `json:"priority"`
}

type AssignConversationInput struct {
	ConversationID int  `json:"conversation_id"`
	AssigneeID     *int `json:"assignee_id,omitempty"`
	TeamID         *int `json:"team_id,omitempty"`
}

type UpdateLabelsInput struct {
	ConversationID int      `json:"conversation_id"`
	Labels         []string `json:"labels"`
}

// RegisterConversationTools registers all conversation-related tools on the MCP server.
func RegisterConversationTools(server *mcp.Server, client *chatwoot.Client) {

	// --- list_conversations ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_conversations",
		Description: "List conversations in Chatwoot. Filter by status: open, resolved, pending, snoozed, all.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListConversationsInput) (*mcp.CallToolResult, any, error) {
		resp, err := client.ListConversations(ctx, input.Status, input.Page)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Total: %d conversations (page %d)\n\n", resp.Data.Meta.AllCount, resp.Data.Meta.Page))
		for _, conv := range resp.Data.Payload {
			assignee := "(unassigned)"
			if conv.Meta.Assignee != nil {
				assignee = conv.Meta.Assignee.Name
			}
			labels := ""
			if len(conv.Labels) > 0 {
				labels = " [" + strings.Join(conv.Labels, ", ") + "]"
			}
			sb.WriteString(fmt.Sprintf("- #%d [%s] %s → %s%s (msgs: %d, unread: %d)\n",
				conv.ID, conv.Status, conv.Meta.Sender.Name, assignee, labels,
				len(conv.Messages), conv.UnreadCount))
		}
		if len(resp.Data.Payload) == 0 {
			sb.WriteString("No conversations found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- get_conversation ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_conversation",
		Description: "Get detailed information about a specific conversation by its ID.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetConversationInput) (*mcp.CallToolResult, any, error) {
		conv, err := client.GetConversation(ctx, input.ConversationID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Conversation #%d\n", conv.ID))
		sb.WriteString(fmt.Sprintf("Status: %s\n", conv.Status))
		sb.WriteString(fmt.Sprintf("Inbox ID: %d\n", conv.InboxID))
		sb.WriteString(fmt.Sprintf("Contact: %s (ID: %d)\n", conv.Meta.Sender.Name, conv.Meta.Sender.ID))
		if conv.Meta.Assignee != nil {
			sb.WriteString(fmt.Sprintf("Assignee: %s (ID: %d)\n", conv.Meta.Assignee.Name, conv.Meta.Assignee.ID))
		} else {
			sb.WriteString("Assignee: (unassigned)\n")
		}
		if conv.Priority != nil {
			sb.WriteString(fmt.Sprintf("Priority: %s\n", *conv.Priority))
		}
		if len(conv.Labels) > 0 {
			sb.WriteString(fmt.Sprintf("Labels: %s\n", strings.Join(conv.Labels, ", ")))
		}
		sb.WriteString(fmt.Sprintf("Messages: %d (unread: %d)\n", len(conv.Messages), conv.UnreadCount))
		if conv.CreatedAt.Valid {
			sb.WriteString(fmt.Sprintf("Created: %s\n", conv.CreatedAt.Format(time.RFC3339)))
		}
		if conv.LastActivityAt.Valid {
			sb.WriteString(fmt.Sprintf("Last activity: %s\n", conv.LastActivityAt.Format(time.RFC3339)))
		}
		return textResult(sb.String()), nil, nil
	})

	// --- create_conversation ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_conversation",
		Description: "Create a new conversation. Requires inbox_id and contact_id. Optionally provide an initial message, status, assignee_id, or team_id.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateConversationInput) (*mcp.CallToolResult, any, error) {
		mbReq := chatwoot.CreateConversationRequest{
			InboxID:   input.InboxID,
			ContactID: input.ContactID,
		}
		if input.Status != "" {
			mbReq.Status = &input.Status
		}
		if input.Message != "" {
			mbReq.Message = &chatwoot.ConversationInitialMessage{Content: input.Message}
		}
		if input.AssigneeID > 0 {
			mbReq.AssigneeID = &input.AssigneeID
		}
		if input.TeamID > 0 {
			mbReq.TeamID = &input.TeamID
		}
		conv, err := client.CreateConversation(ctx, mbReq)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Conversation created! #%d (inbox: %d, status: %s)", conv.ID, conv.InboxID, conv.Status)), nil, nil
	})

	// --- filter_conversations ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "filter_conversations",
		Description: "Filter conversations using advanced criteria. Each filter has attribute_key (status, assignee_id, inbox_id, team_id, label, priority, created_at, last_activity_at, etc.), filter_operator (equal_to, not_equal_to, contains, is_greater_than, is_less_than, days_before, etc.), values (array), and query_operator (AND/OR) to chain filters.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input FilterConversationsInput) (*mcp.CallToolResult, any, error) {
		payload := make([]chatwoot.ConversationFilterPayload, len(input.Payload))
		for i, p := range input.Payload {
			values := make([]any, len(p.Values))
			for j, v := range p.Values {
				if n, err := strconv.Atoi(v); err == nil {
					values[j] = n
				} else {
					values[j] = v
				}
			}
			fp := chatwoot.ConversationFilterPayload{
				AttributeKey:   p.AttributeKey,
				FilterOperator: p.FilterOperator,
				Values:         values,
			}
			if p.QueryOperator != "" {
				fp.QueryOperator = &p.QueryOperator
			}
			payload[i] = fp
		}
		filterReq := chatwoot.ConversationFilterRequest{
			Payload: payload,
		}
		if input.Page > 0 {
			filterReq.Page = &input.Page
		}
		resp, err := client.FilterConversations(ctx, filterReq)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Filter results: %d conversations\n\n", len(resp.Payload)))
		for _, conv := range resp.Payload {
			assignee := "(unassigned)"
			if conv.Meta.Assignee != nil {
				assignee = conv.Meta.Assignee.Name
			}
			sb.WriteString(fmt.Sprintf("- #%d [%s] %s → %s (msgs: %d)\n",
				conv.ID, conv.Status, conv.Meta.Sender.Name, assignee, len(conv.Messages)))
		}
		if len(resp.Payload) == 0 {
			sb.WriteString("No conversations match the filter.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- get_conversation_counts ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_conversation_counts",
		Description: "Get conversation counts grouped by status (open, pending, resolved, snoozed, all).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetConversationMetaInput) (*mcp.CallToolResult, any, error) {
		meta, err := client.GetConversationMeta(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString("Conversation counts:\n")
		sb.WriteString(fmt.Sprintf("  All: %d\n", meta.Meta.AllCount))
		sb.WriteString(fmt.Sprintf("  Mine: %d\n", meta.Meta.MineCount))
		sb.WriteString(fmt.Sprintf("  Assigned: %d\n", meta.Meta.AssignedCount))
		sb.WriteString(fmt.Sprintf("  Unassigned: %d\n", meta.Meta.UnassignedCount))
		return textResult(sb.String()), nil, nil
	})

	// --- update_conversation ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_conversation",
		Description: "Update conversation custom attributes.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input UpdateConversationInput) (*mcp.CallToolResult, any, error) {
		if _, err := client.UpdateConversation(ctx, input.ConversationID, chatwoot.UpdateConversationRequest{
			CustomAttributes: input.CustomAttributes,
		}); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Conversation #%d updated!", input.ConversationID)), nil, nil
	})

	// --- get_messages ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_messages",
		Description: "Get all messages in a conversation. Returns message content, sender, type, and timestamps.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetMessagesInput) (*mcp.CallToolResult, any, error) {
		messages, err := client.GetMessages(ctx, input.ConversationID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Messages in conversation #%d (%d total):\n\n", input.ConversationID, len(messages)))

		maxMessages := 50
		for i, msg := range messages {
			if i >= maxMessages {
				sb.WriteString(fmt.Sprintf("\n... (%d more messages truncated)", len(messages)-maxMessages))
				break
			}
			msgType := messageTypeName(msg.MessageType)
			senderName := "(system)"
			if msg.Sender != nil {
				senderName = msg.Sender.Name
			}
			content := "(no content)"
			if msg.Content != nil && *msg.Content != "" {
				content = *msg.Content
			}
			private := ""
			if msg.Private {
				private = " [private]"
			}
			ts := time.Unix(msg.CreatedAt, 0).Format("2006-01-02 15:04")
			sb.WriteString(fmt.Sprintf("[%s] %s (%s)%s: %s\n", ts, senderName, msgType, private, content))
		}
		if len(messages) == 0 {
			sb.WriteString("No messages found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- send_message ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_message",
		Description: "Send a message to a conversation. message_type: 'outgoing' (to customer) or 'incoming'. Set private=true for internal notes visible only to agents.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SendMessageInput) (*mcp.CallToolResult, any, error) {
		msgType := input.MessageType
		if msgType == "" {
			msgType = "outgoing"
		}
		msg, err := client.SendMessage(ctx, input.ConversationID, chatwoot.SendMessageRequest{
			Content:     input.Content,
			MessageType: msgType,
			Private:     input.Private,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Message sent! (ID: %d, conversation: #%d)", msg.ID, input.ConversationID)), nil, nil
	})

	// --- delete_message ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_message",
		Description: "Delete a specific message from a conversation.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteMessageInput) (*mcp.CallToolResult, any, error) {
		if err := client.DeleteMessage(ctx, input.ConversationID, input.MessageID); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Message #%d deleted from conversation #%d.", input.MessageID, input.ConversationID)), nil, nil
	})

	// --- toggle_conversation_status ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "toggle_conversation_status",
		Description: "Change the status of a conversation. Valid statuses: open, resolved, pending, snoozed.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ToggleStatusInput) (*mcp.CallToolResult, any, error) {
		if err := client.ToggleStatus(ctx, input.ConversationID, input.Status); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Conversation #%d status changed to '%s'!", input.ConversationID, input.Status)), nil, nil
	})

	// --- toggle_conversation_priority ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "toggle_conversation_priority",
		Description: "Set the priority of a conversation. Valid priorities: urgent, high, medium, low, none.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TogglePriorityInput) (*mcp.CallToolResult, any, error) {
		if err := client.TogglePriority(ctx, input.ConversationID, input.Priority); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Conversation #%d priority set to '%s'!", input.ConversationID, input.Priority)), nil, nil
	})

	// --- assign_conversation ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "assign_conversation",
		Description: "Assign a conversation to an agent and/or team. Set assignee_id to null to unassign.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input AssignConversationInput) (*mcp.CallToolResult, any, error) {
		if err := client.AssignConversation(ctx, input.ConversationID, chatwoot.AssignConversationRequest{
			AssigneeID: input.AssigneeID,
			TeamID:     input.TeamID,
		}); err != nil {
			return errorResult(err), nil, nil
		}
		result := fmt.Sprintf("Conversation #%d assignment updated!", input.ConversationID)
		if input.AssigneeID != nil {
			result += fmt.Sprintf(" Agent ID: %d", *input.AssigneeID)
		}
		if input.TeamID != nil {
			result += fmt.Sprintf(" Team ID: %d", *input.TeamID)
		}
		return textResult(result), nil, nil
	})

	// --- update_conversation_labels ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_conversation_labels",
		Description: "Update labels on a conversation. Provide the full list of labels to set (replaces existing).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input UpdateLabelsInput) (*mcp.CallToolResult, any, error) {
		if err := client.UpdateConversationLabels(ctx, input.ConversationID, input.Labels); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Conversation #%d labels updated to: %s", input.ConversationID, strings.Join(input.Labels, ", "))), nil, nil
	})
}

func messageTypeName(t int) string {
	switch t {
	case 0:
		return "incoming"
	case 1:
		return "outgoing"
	case 2:
		return "activity"
	default:
		return fmt.Sprintf("type-%d", t)
	}
}
