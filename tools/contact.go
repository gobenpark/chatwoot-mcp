package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gobenpark/chatwoot-mcp/chatwoot"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- Input types ---

type ListContactsInput struct {
	Page int `json:"page,omitempty"`
}

type GetContactInput struct {
	ContactID int `json:"contact_id"`
}

type SearchContactsInput struct {
	Query string `json:"query"`
}

type CreateContactInput struct {
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	InboxID     int    `json:"inbox_id,omitempty"`
	Identifier  string `json:"identifier,omitempty"`
}

type UpdateContactInput struct {
	ContactID   int    `json:"contact_id"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Identifier  string `json:"identifier,omitempty"`
}

type DeleteContactInput struct {
	ContactID int `json:"contact_id"`
}

type FilterContactsInput struct {
	Payload []FilterPayloadInput `json:"payload"`
	Page    int                  `json:"page,omitempty"`
}

type GetContactConversationsInput struct {
	ContactID int `json:"contact_id"`
}

type MergeContactsInput struct {
	BaseContactID   int `json:"base_contact_id"`
	MergeeContactID int `json:"mergee_contact_id"`
}

type UpdateContactLabelsInput struct {
	ContactID int      `json:"contact_id"`
	Labels    []string `json:"labels"`
}

// RegisterContactTools registers all contact-related tools on the MCP server.
func RegisterContactTools(server *mcp.Server, client *chatwoot.Client) {

	// --- list_contacts ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_contacts",
		Description: "List contacts in the Chatwoot account with pagination.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListContactsInput) (*mcp.CallToolResult, any, error) {
		resp, err := client.ListContacts(ctx, input.Page)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		for _, c := range resp.Payload {
			sb.WriteString(formatContactLine(c))
		}
		if sb.Len() == 0 {
			sb.WriteString("No contacts found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- get_contact ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_contact",
		Description: "Get detailed information about a specific contact by ID.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetContactInput) (*mcp.CallToolResult, any, error) {
		contact, err := client.GetContact(ctx, input.ContactID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(formatContactDetail(contact)), nil, nil
	})

	// --- search_contacts ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_contacts",
		Description: "Search contacts by name, email, phone number, or identifier.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SearchContactsInput) (*mcp.CallToolResult, any, error) {
		resp, err := client.SearchContacts(ctx, input.Query)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Found %d contacts:\n\n", len(resp.Payload)))
		for _, c := range resp.Payload {
			sb.WriteString(formatContactLine(c))
		}
		if len(resp.Payload) == 0 {
			sb.WriteString("No contacts found.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- create_contact ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_contact",
		Description: "Create a new contact in Chatwoot. At minimum, provide a name.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateContactInput) (*mcp.CallToolResult, any, error) {
		contact, err := client.CreateContact(ctx, chatwoot.CreateContactRequest{
			Name:        input.Name,
			Email:       input.Email,
			PhoneNumber: input.PhoneNumber,
			InboxID:     input.InboxID,
			Identifier:  input.Identifier,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Contact created! ID: %d, Name: %s", contact.ID, contact.Name)), nil, nil
	})

	// --- update_contact ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_contact",
		Description: "Update a contact's name, email, phone number, or identifier. Provide only fields you want to change.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input UpdateContactInput) (*mcp.CallToolResult, any, error) {
		mbReq := chatwoot.UpdateContactRequest{}
		if input.Name != "" {
			mbReq.Name = &input.Name
		}
		if input.Email != "" {
			mbReq.Email = &input.Email
		}
		if input.PhoneNumber != "" {
			mbReq.PhoneNumber = &input.PhoneNumber
		}
		if input.Identifier != "" {
			mbReq.Identifier = &input.Identifier
		}
		contact, err := client.UpdateContact(ctx, input.ContactID, mbReq)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Contact #%d updated! Name: %s", contact.ID, contact.Name)), nil, nil
	})

	// --- delete_contact ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_contact",
		Description: "Permanently delete a contact from Chatwoot.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteContactInput) (*mcp.CallToolResult, any, error) {
		if err := client.DeleteContact(ctx, input.ContactID); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Contact #%d deleted.", input.ContactID)), nil, nil
	})

	// --- filter_contacts ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "filter_contacts",
		Description: "Filter contacts using advanced criteria. Each filter: attribute_key (name, email, phone_number, identifier, etc.), filter_operator (equal_to, not_equal_to, contains, etc.), values, query_operator (AND/OR).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input FilterContactsInput) (*mcp.CallToolResult, any, error) {
		payload := make([]chatwoot.ContactFilterPayload, len(input.Payload))
		for i, p := range input.Payload {
			values := make([]any, len(p.Values))
			for j, v := range p.Values {
				if n, err := strconv.Atoi(v); err == nil {
					values[j] = n
				} else {
					values[j] = v
				}
			}
			fp := chatwoot.ContactFilterPayload{
				AttributeKey:   p.AttributeKey,
				FilterOperator: p.FilterOperator,
				Values:         values,
			}
			if p.QueryOperator != "" {
				fp.QueryOperator = &p.QueryOperator
			}
			payload[i] = fp
		}
		filterReq := chatwoot.ContactFilterRequest{Payload: payload}
		if input.Page > 0 {
			filterReq.Page = &input.Page
		}
		resp, err := client.FilterContacts(ctx, filterReq)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Filter results: %d contacts\n\n", len(resp.Payload)))
		for _, c := range resp.Payload {
			sb.WriteString(formatContactLine(c))
		}
		if len(resp.Payload) == 0 {
			sb.WriteString("No contacts match the filter.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- get_contact_conversations ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_contact_conversations",
		Description: "List all conversations for a specific contact.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetContactConversationsInput) (*mcp.CallToolResult, any, error) {
		resp, err := client.GetContactConversations(ctx, input.ContactID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Conversations for contact #%d (%d total):\n\n", input.ContactID, len(resp.Payload)))
		for _, conv := range resp.Payload {
			assignee := "(unassigned)"
			if conv.Meta.Assignee != nil {
				assignee = conv.Meta.Assignee.Name
			}
			sb.WriteString(fmt.Sprintf("- #%d [%s] → %s (msgs: %d)\n",
				conv.ID, conv.Status, assignee, len(conv.Messages)))
		}
		if len(resp.Payload) == 0 {
			sb.WriteString("No conversations found for this contact.")
		}
		return textResult(sb.String()), nil, nil
	})

	// --- merge_contacts ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "merge_contacts",
		Description: "Merge two contacts. The mergee_contact will be merged into base_contact. All conversations from the mergee will be transferred.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MergeContactsInput) (*mcp.CallToolResult, any, error) {
		contact, err := client.MergeContacts(ctx, chatwoot.MergeContactsRequest{
			BaseContactID:   input.BaseContactID,
			MergeeContactID: input.MergeeContactID,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Contacts merged! Resulting contact: #%d (%s)", contact.ID, contact.Name)), nil, nil
	})

	// --- update_contact_labels ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_contact_labels",
		Description: "Update labels on a contact. Provide the full list of labels to set (replaces existing).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input UpdateContactLabelsInput) (*mcp.CallToolResult, any, error) {
		if _, err := client.UpdateContactLabels(ctx, input.ContactID, input.Labels); err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(fmt.Sprintf("Contact #%d labels updated to: %s", input.ContactID, strings.Join(input.Labels, ", "))), nil, nil
	})
}

func formatContactLine(c chatwoot.Contact) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("- [%d] %s", c.ID, c.Name))
	if c.Email != nil && *c.Email != "" {
		sb.WriteString(fmt.Sprintf(" <%s>", *c.Email))
	}
	if c.PhoneNumber != nil && *c.PhoneNumber != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", *c.PhoneNumber))
	}
	sb.WriteString("\n")
	return sb.String()
}

func formatContactDetail(c *chatwoot.Contact) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Contact #%d: %s\n", c.ID, c.Name))
	if c.Email != nil && *c.Email != "" {
		sb.WriteString(fmt.Sprintf("Email: %s\n", *c.Email))
	}
	if c.PhoneNumber != nil && *c.PhoneNumber != "" {
		sb.WriteString(fmt.Sprintf("Phone: %s\n", *c.PhoneNumber))
	}
	if c.Identifier != nil && *c.Identifier != "" {
		sb.WriteString(fmt.Sprintf("Identifier: %s\n", *c.Identifier))
	}
	if c.CreatedAt.Valid {
		sb.WriteString(fmt.Sprintf("Created: %s\n", c.CreatedAt.Format("2006-01-02 15:04")))
	}
	if c.LastActivityAt.Valid {
		sb.WriteString(fmt.Sprintf("Last activity: %s\n", c.LastActivityAt.Format("2006-01-02 15:04")))
	}
	if len(c.CustomAttributes) > 0 {
		sb.WriteString("Custom attributes:\n")
		for k, v := range c.CustomAttributes {
			sb.WriteString(fmt.Sprintf("  %s: %v\n", k, v))
		}
	}
	return sb.String()
}
