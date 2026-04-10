package chatwoot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Client communicates with the Chatwoot REST API.
type Client struct {
	baseURL    string
	apiToken   string
	accountID  int
	httpClient *http.Client
}

// NewClientFromEnv creates a Client from CHATWOOT_URL, CHATWOOT_API_TOKEN, and CHATWOOT_ACCOUNT_ID env vars.
func NewClientFromEnv() (*Client, error) {
	baseURL := strings.TrimRight(os.Getenv("CHATWOOT_URL"), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("CHATWOOT_URL environment variable is required")
	}
	apiToken := os.Getenv("CHATWOOT_API_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("CHATWOOT_API_TOKEN environment variable is required")
	}
	accountIDStr := os.Getenv("CHATWOOT_ACCOUNT_ID")
	if accountIDStr == "" {
		return nil, fmt.Errorf("CHATWOOT_ACCOUNT_ID environment variable is required")
	}
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		return nil, fmt.Errorf("CHATWOOT_ACCOUNT_ID must be a valid integer: %w", err)
	}
	return &Client{
		baseURL:   baseURL,
		apiToken:  apiToken,
		accountID: accountID,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

func (c *Client) accountPath(path string) string {
	return fmt.Sprintf("/api/v1/accounts/%d%s", c.accountID, path)
}

func (c *Client) accountPathV2(path string) string {
	return fmt.Sprintf("/api/v2/accounts/%d%s", c.accountID, path)
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("api_access_token", c.apiToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("chatwoot API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Conversations
// ---------------------------------------------------------------------------

// ListConversations returns conversations with optional status filter and pagination.
func (c *Client) ListConversations(ctx context.Context, status, assigneeType, q string, inboxID, teamID int, labels []string, page int) (*ConversationListResponse, error) {
	params := url.Values{}
	if status != "" {
		params.Set("status", status)
	}
	if assigneeType != "" {
		params.Set("assignee_type", assigneeType)
	}
	if q != "" {
		params.Set("q", q)
	}
	if inboxID > 0 {
		params.Set("inbox_id", strconv.Itoa(inboxID))
	}
	if teamID > 0 {
		params.Set("team_id", strconv.Itoa(teamID))
	}
	for _, l := range labels {
		params.Add("labels[]", l)
	}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	path := c.accountPath("/conversations")
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	var resp ConversationListResponse
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetConversation returns a single conversation by ID.
func (c *Client) GetConversation(ctx context.Context, conversationID int) (*Conversation, error) {
	var conv Conversation
	path := c.accountPath(fmt.Sprintf("/conversations/%d", conversationID))
	if err := c.do(ctx, http.MethodGet, path, nil, &conv); err != nil {
		return nil, err
	}
	return &conv, nil
}

// CreateConversation creates a new conversation.
func (c *Client) CreateConversation(ctx context.Context, req CreateConversationRequest) (*Conversation, error) {
	var conv Conversation
	path := c.accountPath("/conversations")
	if err := c.do(ctx, http.MethodPost, path, req, &conv); err != nil {
		return nil, err
	}
	return &conv, nil
}

// FilterConversations filters conversations by the given criteria.
func (c *Client) FilterConversations(ctx context.Context, req ConversationFilterRequest) (*ConversationFilterResponse, error) {
	var resp ConversationFilterResponse
	path := c.accountPath("/conversations/filter")
	if req.Page != nil && *req.Page > 0 {
		path += "?page=" + strconv.Itoa(*req.Page)
	}
	if err := c.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetConversationMeta returns conversation counts per status.
func (c *Client) GetConversationMeta(ctx context.Context) (*ConversationMetaResponse, error) {
	var resp ConversationMetaResponse
	path := c.accountPath("/conversations/meta")
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateConversation updates a conversation's custom attributes.
func (c *Client) UpdateConversation(ctx context.Context, conversationID int, req UpdateConversationRequest) (*Conversation, error) {
	var conv Conversation
	path := c.accountPath(fmt.Sprintf("/conversations/%d", conversationID))
	if err := c.do(ctx, http.MethodPatch, path, req, &conv); err != nil {
		return nil, err
	}
	return &conv, nil
}

// TogglePriority sets the priority of a conversation.
func (c *Client) TogglePriority(ctx context.Context, conversationID int, priority string) error {
	path := c.accountPath(fmt.Sprintf("/conversations/%d/toggle_priority", conversationID))
	payload := TogglePriorityRequest{Priority: priority}
	return c.do(ctx, http.MethodPost, path, payload, nil)
}

// ToggleStatus toggles the status of a conversation (open, resolved, pending).
func (c *Client) ToggleStatus(ctx context.Context, conversationID int, status string) error {
	path := c.accountPath(fmt.Sprintf("/conversations/%d/toggle_status", conversationID))
	payload := ToggleStatusRequest{Status: status}
	return c.do(ctx, http.MethodPost, path, payload, nil)
}

// AssignConversation assigns a conversation to an agent and/or team.
func (c *Client) AssignConversation(ctx context.Context, conversationID int, req AssignConversationRequest) error {
	path := c.accountPath(fmt.Sprintf("/conversations/%d/assignments", conversationID))
	return c.do(ctx, http.MethodPost, path, req, nil)
}

// UpdateConversationLabels updates labels on a conversation.
func (c *Client) UpdateConversationLabels(ctx context.Context, conversationID int, labels []string) error {
	path := c.accountPath(fmt.Sprintf("/conversations/%d/labels", conversationID))
	payload := ConversationLabelsRequest{Labels: labels}
	return c.do(ctx, http.MethodPost, path, payload, nil)
}

// GetConversationLabels gets labels for a conversation.
func (c *Client) GetConversationLabels(ctx context.Context, conversationID int) ([]string, error) {
	var resp struct {
		Payload []string `json:"payload"`
	}
	path := c.accountPath(fmt.Sprintf("/conversations/%d/labels", conversationID))
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------

// GetMessages returns messages for a conversation.
func (c *Client) GetMessages(ctx context.Context, conversationID int) ([]Message, error) {
	var resp MessageListResponse
	path := c.accountPath(fmt.Sprintf("/conversations/%d/messages", conversationID))
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// SendMessage sends a message to a conversation.
func (c *Client) SendMessage(ctx context.Context, conversationID int, req SendMessageRequest) (*Message, error) {
	var msg Message
	path := c.accountPath(fmt.Sprintf("/conversations/%d/messages", conversationID))
	if err := c.do(ctx, http.MethodPost, path, req, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// DeleteMessage deletes a message from a conversation.
func (c *Client) DeleteMessage(ctx context.Context, conversationID, messageID int) error {
	path := c.accountPath(fmt.Sprintf("/conversations/%d/messages/%d", conversationID, messageID))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// ---------------------------------------------------------------------------
// Contacts
// ---------------------------------------------------------------------------

// ListContacts returns contacts with pagination.
func (c *Client) ListContacts(ctx context.Context, page int) (*ContactListResponse, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	path := c.accountPath("/contacts")
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	var resp ContactListResponse
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetContact returns a single contact by ID.
func (c *Client) GetContact(ctx context.Context, contactID int) (*Contact, error) {
	var resp struct {
		Payload Contact `json:"payload"`
	}
	path := c.accountPath(fmt.Sprintf("/contacts/%d", contactID))
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Payload, nil
}

// SearchContacts searches contacts by query string.
func (c *Client) SearchContacts(ctx context.Context, query string) (*ContactSearchResponse, error) {
	var resp ContactSearchResponse
	path := c.accountPath("/contacts/search?q=" + url.QueryEscape(query))
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateContact creates a new contact.
func (c *Client) CreateContact(ctx context.Context, req CreateContactRequest) (*Contact, error) {
	var resp struct {
		Payload struct {
			Contact Contact `json:"contact"`
		} `json:"payload"`
	}
	path := c.accountPath("/contacts")
	if err := c.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp.Payload.Contact, nil
}

// UpdateContact updates an existing contact.
func (c *Client) UpdateContact(ctx context.Context, contactID int, req UpdateContactRequest) (*Contact, error) {
	var resp struct {
		Payload Contact `json:"payload"`
	}
	path := c.accountPath(fmt.Sprintf("/contacts/%d", contactID))
	if err := c.do(ctx, http.MethodPut, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp.Payload, nil
}

// DeleteContact deletes a contact by ID.
func (c *Client) DeleteContact(ctx context.Context, contactID int) error {
	path := c.accountPath(fmt.Sprintf("/contacts/%d", contactID))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// FilterContacts filters contacts by the given criteria.
func (c *Client) FilterContacts(ctx context.Context, req ContactFilterRequest) (*ContactListResponse, error) {
	var resp ContactListResponse
	path := c.accountPath("/contacts/filter")
	if req.Page != nil && *req.Page > 0 {
		path += "?page=" + strconv.Itoa(*req.Page)
	}
	if err := c.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetContactConversations returns conversations for a given contact.
func (c *Client) GetContactConversations(ctx context.Context, contactID int) (*ContactConversationsResponse, error) {
	var resp ContactConversationsResponse
	path := c.accountPath(fmt.Sprintf("/contacts/%d/conversations", contactID))
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MergeContacts merges two contacts together.
func (c *Client) MergeContacts(ctx context.Context, req MergeContactsRequest) (*Contact, error) {
	var resp struct {
		Payload Contact `json:"payload"`
	}
	path := c.accountPath("/actions/contact_merge")
	if err := c.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp.Payload, nil
}

// GetContactLabels returns labels for a contact.
func (c *Client) GetContactLabels(ctx context.Context, contactID int) ([]string, error) {
	var resp struct {
		Payload []string `json:"payload"`
	}
	path := c.accountPath(fmt.Sprintf("/contacts/%d/labels", contactID))
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// UpdateContactLabels updates labels on a contact.
func (c *Client) UpdateContactLabels(ctx context.Context, contactID int, labels []string) ([]string, error) {
	var resp struct {
		Payload []string `json:"payload"`
	}
	path := c.accountPath(fmt.Sprintf("/contacts/%d/labels", contactID))
	payload := struct {
		Labels []string `json:"labels"`
	}{Labels: labels}
	if err := c.do(ctx, http.MethodPost, path, payload, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// ---------------------------------------------------------------------------
// Inboxes
// ---------------------------------------------------------------------------

// ListInboxes returns all inboxes for the account.
func (c *Client) ListInboxes(ctx context.Context) ([]Inbox, error) {
	var resp InboxListResponse
	path := c.accountPath("/inboxes")
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// ---------------------------------------------------------------------------
// Agents
// ---------------------------------------------------------------------------

// ListAgents returns all agents for the account.
func (c *Client) ListAgents(ctx context.Context) ([]Agent, error) {
	var agents []Agent
	path := c.accountPath("/agents")
	if err := c.do(ctx, http.MethodGet, path, nil, &agents); err != nil {
		return nil, err
	}
	return agents, nil
}

// ---------------------------------------------------------------------------
// Labels
// ---------------------------------------------------------------------------

// ListLabels returns all labels for the account.
func (c *Client) ListLabels(ctx context.Context) ([]Label, error) {
	var resp struct {
		Payload []Label `json:"payload"`
	}
	path := c.accountPath("/labels")
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// ---------------------------------------------------------------------------
// Teams
// ---------------------------------------------------------------------------

// ListTeams returns all teams for the account.
func (c *Client) ListTeams(ctx context.Context) ([]Team, error) {
	var teams []Team
	path := c.accountPath("/teams")
	if err := c.do(ctx, http.MethodGet, path, nil, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

// GetTeam returns a single team by ID.
func (c *Client) GetTeam(ctx context.Context, teamID int) (*Team, error) {
	var team Team
	path := c.accountPath(fmt.Sprintf("/teams/%d", teamID))
	if err := c.do(ctx, http.MethodGet, path, nil, &team); err != nil {
		return nil, err
	}
	return &team, nil
}

// CreateTeam creates a new team.
func (c *Client) CreateTeam(ctx context.Context, name, description string) (*Team, error) {
	var team Team
	path := c.accountPath("/teams")
	payload := struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}{Name: name, Description: description}
	if err := c.do(ctx, http.MethodPost, path, payload, &team); err != nil {
		return nil, err
	}
	return &team, nil
}

// UpdateTeam updates a team.
func (c *Client) UpdateTeam(ctx context.Context, teamID int, name, description string) (*Team, error) {
	var team Team
	path := c.accountPath(fmt.Sprintf("/teams/%d", teamID))
	payload := struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	}{Name: name, Description: description}
	if err := c.do(ctx, http.MethodPatch, path, payload, &team); err != nil {
		return nil, err
	}
	return &team, nil
}

// DeleteTeam deletes a team by ID.
func (c *Client) DeleteTeam(ctx context.Context, teamID int) error {
	path := c.accountPath(fmt.Sprintf("/teams/%d", teamID))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// ListTeamMembers returns the members of a team.
func (c *Client) ListTeamMembers(ctx context.Context, teamID int) ([]TeamMember, error) {
	var members []TeamMember
	path := c.accountPath(fmt.Sprintf("/teams/%d/team_members", teamID))
	if err := c.do(ctx, http.MethodGet, path, nil, &members); err != nil {
		return nil, err
	}
	return members, nil
}

// AddTeamMembers adds members to a team.
func (c *Client) AddTeamMembers(ctx context.Context, teamID int, userIDs []int) ([]TeamMember, error) {
	var members []TeamMember
	path := c.accountPath(fmt.Sprintf("/teams/%d/team_members", teamID))
	payload := AddTeamMemberRequest{UserIDs: userIDs}
	if err := c.do(ctx, http.MethodPost, path, payload, &members); err != nil {
		return nil, err
	}
	return members, nil
}

// RemoveTeamMembers removes members from a team.
func (c *Client) RemoveTeamMembers(ctx context.Context, teamID int, userIDs []int) error {
	path := c.accountPath(fmt.Sprintf("/teams/%d/team_members", teamID))
	payload := AddTeamMemberRequest{UserIDs: userIDs}
	return c.do(ctx, http.MethodDelete, path, payload, nil)
}

// ---------------------------------------------------------------------------
// Canned Responses
// ---------------------------------------------------------------------------

// ListCannedResponses returns all canned responses for the account.
func (c *Client) ListCannedResponses(ctx context.Context) ([]CannedResponse, error) {
	var responses []CannedResponse
	path := c.accountPath("/canned_responses")
	if err := c.do(ctx, http.MethodGet, path, nil, &responses); err != nil {
		return nil, err
	}
	return responses, nil
}

// CreateCannedResponse creates a new canned response.
func (c *Client) CreateCannedResponse(ctx context.Context, req CreateCannedResponseRequest) (*CannedResponse, error) {
	var resp CannedResponse
	path := c.accountPath("/canned_responses")
	if err := c.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCannedResponse updates a canned response by ID.
func (c *Client) UpdateCannedResponse(ctx context.Context, id int, req UpdateCannedResponseRequest) (*CannedResponse, error) {
	var resp CannedResponse
	path := c.accountPath(fmt.Sprintf("/canned_responses/%d", id))
	if err := c.do(ctx, http.MethodPatch, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCannedResponse deletes a canned response by ID.
func (c *Client) DeleteCannedResponse(ctx context.Context, id int) error {
	path := c.accountPath(fmt.Sprintf("/canned_responses/%d", id))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// ---------------------------------------------------------------------------
// Custom Attribute Definitions
// ---------------------------------------------------------------------------

// ListCustomAttributes returns all custom attribute definitions for the account.
func (c *Client) ListCustomAttributes(ctx context.Context, attributeModel int) ([]CustomAttributeDefinition, error) {
	var attrs []CustomAttributeDefinition
	path := c.accountPath("/custom_attribute_definitions")
	if attributeModel >= 0 {
		path += "?attribute_model=" + strconv.Itoa(attributeModel)
	}
	if err := c.do(ctx, http.MethodGet, path, nil, &attrs); err != nil {
		return nil, err
	}
	return attrs, nil
}

// CreateCustomAttribute creates a new custom attribute definition.
func (c *Client) CreateCustomAttribute(ctx context.Context, req CreateCustomAttributeRequest) (*CustomAttributeDefinition, error) {
	var attr CustomAttributeDefinition
	path := c.accountPath("/custom_attribute_definitions")
	if err := c.do(ctx, http.MethodPost, path, req, &attr); err != nil {
		return nil, err
	}
	return &attr, nil
}

// UpdateCustomAttribute updates a custom attribute definition by ID.
func (c *Client) UpdateCustomAttribute(ctx context.Context, id int, req CreateCustomAttributeRequest) (*CustomAttributeDefinition, error) {
	var attr CustomAttributeDefinition
	path := c.accountPath(fmt.Sprintf("/custom_attribute_definitions/%d", id))
	if err := c.do(ctx, http.MethodPatch, path, req, &attr); err != nil {
		return nil, err
	}
	return &attr, nil
}

// DeleteCustomAttribute deletes a custom attribute definition by ID.
func (c *Client) DeleteCustomAttribute(ctx context.Context, id int) error {
	path := c.accountPath(fmt.Sprintf("/custom_attribute_definitions/%d", id))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// ---------------------------------------------------------------------------
// Custom Filters
// ---------------------------------------------------------------------------

// ListCustomFilters returns all custom filters for the account.
func (c *Client) ListCustomFilters(ctx context.Context, filterType string) ([]CustomFilter, error) {
	var filters []CustomFilter
	path := c.accountPath("/custom_filters")
	if filterType != "" {
		path += "?filter_type=" + url.QueryEscape(filterType)
	}
	if err := c.do(ctx, http.MethodGet, path, nil, &filters); err != nil {
		return nil, err
	}
	return filters, nil
}

// CreateCustomFilter creates a new custom filter.
func (c *Client) CreateCustomFilter(ctx context.Context, req CreateCustomFilterRequest) (*CustomFilter, error) {
	var filter CustomFilter
	path := c.accountPath("/custom_filters")
	if err := c.do(ctx, http.MethodPost, path, req, &filter); err != nil {
		return nil, err
	}
	return &filter, nil
}

// UpdateCustomFilter updates a custom filter by ID.
func (c *Client) UpdateCustomFilter(ctx context.Context, id int, req CreateCustomFilterRequest) (*CustomFilter, error) {
	var filter CustomFilter
	path := c.accountPath(fmt.Sprintf("/custom_filters/%d", id))
	if err := c.do(ctx, http.MethodPatch, path, req, &filter); err != nil {
		return nil, err
	}
	return &filter, nil
}

// DeleteCustomFilter deletes a custom filter by ID.
func (c *Client) DeleteCustomFilter(ctx context.Context, id int) error {
	path := c.accountPath(fmt.Sprintf("/custom_filters/%d", id))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// ---------------------------------------------------------------------------
// Automation Rules
// ---------------------------------------------------------------------------

// ListAutomationRules returns all automation rules for the account.
func (c *Client) ListAutomationRules(ctx context.Context) ([]AutomationRule, error) {
	var resp struct {
		Payload []AutomationRule `json:"payload"`
	}
	path := c.accountPath("/automation_rules")
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// CreateAutomationRule creates a new automation rule.
func (c *Client) CreateAutomationRule(ctx context.Context, req CreateAutomationRuleRequest) (*AutomationRule, error) {
	var rule AutomationRule
	path := c.accountPath("/automation_rules")
	if err := c.do(ctx, http.MethodPost, path, req, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// UpdateAutomationRule updates an automation rule by ID.
func (c *Client) UpdateAutomationRule(ctx context.Context, id int, req CreateAutomationRuleRequest) (*AutomationRule, error) {
	var rule AutomationRule
	path := c.accountPath(fmt.Sprintf("/automation_rules/%d", id))
	if err := c.do(ctx, http.MethodPatch, path, req, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// DeleteAutomationRule deletes an automation rule by ID.
func (c *Client) DeleteAutomationRule(ctx context.Context, id int) error {
	path := c.accountPath(fmt.Sprintf("/automation_rules/%d", id))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// ---------------------------------------------------------------------------
// Webhooks
// ---------------------------------------------------------------------------

// ListWebhooks returns all webhooks for the account.
func (c *Client) ListWebhooks(ctx context.Context) ([]Webhook, error) {
	var resp struct {
		Payload []Webhook `json:"payload"`
	}
	path := c.accountPath("/webhooks")
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// CreateWebhook creates a new webhook.
func (c *Client) CreateWebhook(ctx context.Context, req CreateWebhookRequest) (*Webhook, error) {
	var webhook Webhook
	path := c.accountPath("/webhooks")
	if err := c.do(ctx, http.MethodPost, path, req, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// UpdateWebhook updates a webhook by ID.
func (c *Client) UpdateWebhook(ctx context.Context, id int, req UpdateWebhookRequest) (*Webhook, error) {
	var webhook Webhook
	path := c.accountPath(fmt.Sprintf("/webhooks/%d", id))
	if err := c.do(ctx, http.MethodPatch, path, req, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// DeleteWebhook deletes a webhook by ID.
func (c *Client) DeleteWebhook(ctx context.Context, id int) error {
	path := c.accountPath(fmt.Sprintf("/webhooks/%d", id))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// ---------------------------------------------------------------------------
// Reports (v2 API)
// ---------------------------------------------------------------------------

// GetReportsSummary returns aggregate report metrics for the given time range.
// reportType is typically "account", "agent", "team", or "inbox".
func (c *Client) GetReportsSummary(ctx context.Context, since, until int64, reportType string) (*ReportSummary, error) {
	params := url.Values{}
	params.Set("since", strconv.FormatInt(since, 10))
	params.Set("until", strconv.FormatInt(until, 10))
	if reportType != "" {
		params.Set("type", reportType)
	}
	var summary ReportSummary
	path := c.accountPathV2("/reports/summary") + "?" + params.Encode()
	if err := c.do(ctx, http.MethodGet, path, nil, &summary); err != nil {
		return nil, err
	}
	return &summary, nil
}

// GetAgentSummary returns per-agent report metrics for the given time range.
func (c *Client) GetAgentSummary(ctx context.Context, since, until int64) ([]SummaryReportEntry, error) {
	params := url.Values{}
	params.Set("since", strconv.FormatInt(since, 10))
	params.Set("until", strconv.FormatInt(until, 10))
	var agents []SummaryReportEntry
	path := c.accountPathV2("/summary_reports/agent") + "?" + params.Encode()
	if err := c.do(ctx, http.MethodGet, path, nil, &agents); err != nil {
		return nil, err
	}
	return agents, nil
}

// GetTeamSummary returns per-team report metrics for the given time range.
func (c *Client) GetTeamSummary(ctx context.Context, since, until int64) ([]SummaryReportEntry, error) {
	params := url.Values{}
	params.Set("since", strconv.FormatInt(since, 10))
	params.Set("until", strconv.FormatInt(until, 10))
	var teams []SummaryReportEntry
	path := c.accountPathV2("/summary_reports/team") + "?" + params.Encode()
	if err := c.do(ctx, http.MethodGet, path, nil, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

// GetInboxSummary returns per-inbox report metrics for the given time range.
func (c *Client) GetInboxSummary(ctx context.Context, since, until int64) ([]SummaryReportEntry, error) {
	params := url.Values{}
	params.Set("since", strconv.FormatInt(since, 10))
	params.Set("until", strconv.FormatInt(until, 10))
	var inboxes []SummaryReportEntry
	path := c.accountPathV2("/summary_reports/inbox") + "?" + params.Encode()
	if err := c.do(ctx, http.MethodGet, path, nil, &inboxes); err != nil {
		return nil, err
	}
	return inboxes, nil
}

// GetChannelSummary returns per-channel report metrics for the given time range.
func (c *Client) GetChannelSummary(ctx context.Context, since, until int64) ([]ChannelSummary, error) {
	params := url.Values{}
	params.Set("since", strconv.FormatInt(since, 10))
	params.Set("until", strconv.FormatInt(until, 10))
	var raw json.RawMessage
	path := c.accountPathV2("/summary_reports/channel") + "?" + params.Encode()
	if err := c.do(ctx, http.MethodGet, path, nil, &raw); err != nil {
		return nil, err
	}
	// API may return {} (empty object) or [] (array)
	if len(raw) == 0 || string(raw) == "{}" || string(raw) == "null" {
		return nil, nil
	}
	var channels []ChannelSummary
	if err := json.Unmarshal(raw, &channels); err != nil {
		return nil, nil
	}
	return channels, nil
}

// ---------------------------------------------------------------------------
// Inbox Members
// ---------------------------------------------------------------------------

// ListInboxMembers returns the members (agents) of an inbox.
func (c *Client) ListInboxMembers(ctx context.Context, inboxID int) ([]InboxMember, error) {
	var resp struct {
		Payload []InboxMember `json:"payload"`
	}
	path := c.accountPath(fmt.Sprintf("/inbox_members/%d", inboxID))
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// AddInboxMembers adds agents to an inbox.
func (c *Client) AddInboxMembers(ctx context.Context, inboxID int, userIDs []int) ([]InboxMember, error) {
	var resp struct {
		Payload []InboxMember `json:"payload"`
	}
	path := c.accountPath("/inbox_members")
	payload := struct {
		InboxID int   `json:"inbox_id"`
		UserIDs []int `json:"user_ids"`
	}{InboxID: inboxID, UserIDs: userIDs}
	if err := c.do(ctx, http.MethodPost, path, payload, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// ---------------------------------------------------------------------------
// Profile
// ---------------------------------------------------------------------------

// GetProfile returns the authenticated user's profile.
func (c *Client) GetProfile(ctx context.Context) (*Profile, error) {
	var profile Profile
	path := "/api/v1/profile"
	if err := c.do(ctx, http.MethodGet, path, nil, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// ---------------------------------------------------------------------------
// Account Management
// ---------------------------------------------------------------------------

// GetAccount returns the account details.
func (c *Client) GetAccount(ctx context.Context) (*Account, error) {
	var account Account
	path := fmt.Sprintf("/api/v1/accounts/%d", c.accountID)
	if err := c.do(ctx, http.MethodGet, path, nil, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// UpdateAccount updates account details.
func (c *Client) UpdateAccount(ctx context.Context, req UpdateAccountRequest) (*Account, error) {
	var account Account
	path := fmt.Sprintf("/api/v1/accounts/%d", c.accountID)
	if err := c.do(ctx, http.MethodPatch, path, req, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// ---------------------------------------------------------------------------
// Audit Logs
// ---------------------------------------------------------------------------

// ListAuditLogs returns audit log entries with pagination.
func (c *Client) ListAuditLogs(ctx context.Context, page int) (*AuditLogListResponse, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	path := c.accountPath("/audit_logs")
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	var resp AuditLogListResponse
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ---------------------------------------------------------------------------
// Agent Bots
// ---------------------------------------------------------------------------

// ListAgentBots returns all agent bots for the account.
func (c *Client) ListAgentBots(ctx context.Context) ([]AgentBot, error) {
	var bots []AgentBot
	path := c.accountPath("/agent_bots")
	if err := c.do(ctx, http.MethodGet, path, nil, &bots); err != nil {
		return nil, err
	}
	return bots, nil
}

// GetAgentBot returns a single agent bot by ID.
func (c *Client) GetAgentBot(ctx context.Context, botID int) (*AgentBot, error) {
	var bot AgentBot
	path := c.accountPath(fmt.Sprintf("/agent_bots/%d", botID))
	if err := c.do(ctx, http.MethodGet, path, nil, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

// CreateAgentBot creates a new agent bot.
func (c *Client) CreateAgentBot(ctx context.Context, req CreateAgentBotRequest) (*AgentBot, error) {
	var bot AgentBot
	path := c.accountPath("/agent_bots")
	if err := c.do(ctx, http.MethodPost, path, req, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

// UpdateAgentBot updates an agent bot.
func (c *Client) UpdateAgentBot(ctx context.Context, botID int, req UpdateAgentBotRequest) (*AgentBot, error) {
	var bot AgentBot
	path := c.accountPath(fmt.Sprintf("/agent_bots/%d", botID))
	if err := c.do(ctx, http.MethodPatch, path, req, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

// DeleteAgentBot deletes an agent bot by ID.
func (c *Client) DeleteAgentBot(ctx context.Context, botID int) error {
	path := c.accountPath(fmt.Sprintf("/agent_bots/%d", botID))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// ---------------------------------------------------------------------------
// Integrations
// ---------------------------------------------------------------------------

// ListIntegrationApps returns all integrations available for the account.
func (c *Client) ListIntegrationApps(ctx context.Context) ([]IntegrationApp, error) {
	var resp struct {
		Payload []IntegrationApp `json:"payload"`
	}
	path := c.accountPath("/integrations/apps")
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// CreateIntegrationHook creates a new integration hook.
func (c *Client) CreateIntegrationHook(ctx context.Context, req CreateIntegrationHookRequest) (*IntegrationHook, error) {
	var hook IntegrationHook
	path := c.accountPath("/integrations/hooks")
	if err := c.do(ctx, http.MethodPost, path, req, &hook); err != nil {
		return nil, err
	}
	return &hook, nil
}

// UpdateIntegrationHook updates an integration hook.
func (c *Client) UpdateIntegrationHook(ctx context.Context, hookID int, req UpdateIntegrationHookRequest) (*IntegrationHook, error) {
	var hook IntegrationHook
	path := c.accountPath(fmt.Sprintf("/integrations/hooks/%d", hookID))
	if err := c.do(ctx, http.MethodPatch, path, req, &hook); err != nil {
		return nil, err
	}
	return &hook, nil
}

// DeleteIntegrationHook deletes an integration hook.
func (c *Client) DeleteIntegrationHook(ctx context.Context, hookID int) error {
	path := c.accountPath(fmt.Sprintf("/integrations/hooks/%d", hookID))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// ---------------------------------------------------------------------------
// Help Center
// ---------------------------------------------------------------------------

// ListPortals returns all help center portals.
func (c *Client) ListPortals(ctx context.Context) ([]Portal, error) {
	var resp struct {
		Payload []Portal `json:"payload"`
	}
	path := c.accountPath("/portals")
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// ListArticles returns articles for a given portal.
func (c *Client) ListArticles(ctx context.Context, portalID int) ([]Article, error) {
	var resp struct {
		Payload []Article `json:"payload"`
	}
	path := c.accountPath(fmt.Sprintf("/portals/%d/articles", portalID))
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// CreateArticle creates a new article in a portal.
func (c *Client) CreateArticle(ctx context.Context, portalID int, req CreateArticleRequest) (*Article, error) {
	var article Article
	path := c.accountPath(fmt.Sprintf("/portals/%d/articles", portalID))
	if err := c.do(ctx, http.MethodPost, path, req, &article); err != nil {
		return nil, err
	}
	return &article, nil
}

// UpdateArticle updates an article in a portal.
func (c *Client) UpdateArticle(ctx context.Context, portalID, articleID int, req UpdateArticleRequest) (*Article, error) {
	var article Article
	path := c.accountPath(fmt.Sprintf("/portals/%d/articles/%d", portalID, articleID))
	if err := c.do(ctx, http.MethodPut, path, req, &article); err != nil {
		return nil, err
	}
	return &article, nil
}

// DeleteArticle deletes an article from a portal.
func (c *Client) DeleteArticle(ctx context.Context, portalID, articleID int) error {
	path := c.accountPath(fmt.Sprintf("/portals/%d/articles/%d", portalID, articleID))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// UpdatePortal updates a portal.
func (c *Client) UpdatePortal(ctx context.Context, portalID int, req UpdatePortalRequest) (*Portal, error) {
	var portal Portal
	path := c.accountPath(fmt.Sprintf("/portals/%d", portalID))
	if err := c.do(ctx, http.MethodPatch, path, req, &portal); err != nil {
		return nil, err
	}
	return &portal, nil
}

// ListCategories returns categories for a portal.
func (c *Client) ListCategories(ctx context.Context, portalID int) ([]Category, error) {
	var resp struct {
		Payload []Category `json:"payload"`
	}
	path := c.accountPath(fmt.Sprintf("/portals/%d/categories", portalID))
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// CreateCategory creates a new category in a portal.
func (c *Client) CreateCategory(ctx context.Context, portalID int, req CreateCategoryRequest) (*Category, error) {
	var category Category
	path := c.accountPath(fmt.Sprintf("/portals/%d/categories", portalID))
	if err := c.do(ctx, http.MethodPost, path, req, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory updates a category in a portal.
func (c *Client) UpdateCategory(ctx context.Context, portalID, categoryID int, req UpdateCategoryRequest) (*Category, error) {
	var category Category
	path := c.accountPath(fmt.Sprintf("/portals/%d/categories/%d", portalID, categoryID))
	if err := c.do(ctx, http.MethodPatch, path, req, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

// DeleteCategory deletes a category from a portal.
func (c *Client) DeleteCategory(ctx context.Context, portalID, categoryID int) error {
	path := c.accountPath(fmt.Sprintf("/portals/%d/categories/%d", portalID, categoryID))
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// ---------------------------------------------------------------------------
// Notifications
// ---------------------------------------------------------------------------

// ListNotifications returns notifications with pagination.
func (c *Client) ListNotifications(ctx context.Context, page int) ([]Notification, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	var resp struct {
		Data struct {
			Meta    PaginationMeta `json:"meta"`
			Payload []Notification `json:"payload"`
		} `json:"data"`
	}
	path := c.accountPath("/notifications")
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data.Payload, nil
}
