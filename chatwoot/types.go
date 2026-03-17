package chatwoot

import (
	"encoding/json"
	"time"
)

// ---------------------------------------------------------------------------
// Conversation
// ---------------------------------------------------------------------------

// Conversation represents a Chatwoot conversation.
type Conversation struct {
	ID                  int              `json:"id"`
	AccountID           int              `json:"account_id"`
	InboxID             int              `json:"inbox_id"`
	Status              string           `json:"status"`
	MessagesCount       int              `json:"messages_count"`
	UnreadMessagesCount int              `json:"unread_messages_count"`
	AssigneeID          *int             `json:"assignee_id"`
	Meta                ConversationMeta `json:"meta"`
	Labels              []string         `json:"labels"`
	Priority            *string          `json:"priority"`
	CustomAttributes    map[string]any   `json:"custom_attributes"`
	LastActivityAt      time.Time        `json:"last_activity_at"`
	CreatedAt           time.Time        `json:"created_at"`
}

// ConversationMeta holds sender and assignee info.
type ConversationMeta struct {
	Sender   MetaContact `json:"sender"`
	Assignee *MetaAgent  `json:"assignee"`
}

// MetaContact is a lightweight contact reference.
type MetaContact struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MetaAgent is a lightweight agent reference.
type MetaAgent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ConversationListResponse is the response from list conversations.
type ConversationListResponse struct {
	Data DataPayload `json:"data"`
}

// DataPayload wraps the conversation list with metadata.
type DataPayload struct {
	Meta    PaginationMeta `json:"meta"`
	Payload []Conversation `json:"payload"`
}

// PaginationMeta holds pagination info.
type PaginationMeta struct {
	AllCount  int `json:"all_count"`
	Page      int `json:"page"`
	TotalPage int `json:"total_pages"`
}

// CreateConversationRequest is the payload for creating a new conversation.
type CreateConversationRequest struct {
	SourceID   *string          `json:"source_id,omitempty"`
	InboxID    int              `json:"inbox_id"`
	ContactID  int              `json:"contact_id"`
	Message    *ConversationInitialMessage `json:"message,omitempty"`
	Status     *string          `json:"status,omitempty"`
	AssigneeID *int             `json:"assignee_id,omitempty"`
	TeamID     *int             `json:"team_id,omitempty"`
}

// ConversationInitialMessage is the initial message when creating a conversation.
type ConversationInitialMessage struct {
	Content string `json:"content"`
}

// ConversationFilterRequest is the payload for filtering conversations.
type ConversationFilterRequest struct {
	Payload []ConversationFilterPayload `json:"payload"`
	Page    *int                        `json:"page,omitempty"`
}

// ConversationFilterPayload describes a single filter condition.
type ConversationFilterPayload struct {
	AttributeKey   string   `json:"attribute_key"`
	FilterOperator string   `json:"filter_operator"`
	Values         []string `json:"values"`
	QueryOperator  *string  `json:"query_operator,omitempty"`
}

// ConversationMetaResponse is the response containing conversation counts per status.
type ConversationMetaResponse struct {
	Meta ConversationStatusCounts `json:"meta"`
}

// ConversationStatusCounts holds the count of conversations per status.
type ConversationStatusCounts struct {
	Open       int `json:"open"`
	Resolved   int `json:"resolved"`
	Pending    int `json:"pending"`
	Snoozed    int `json:"snoozed"`
	AllCount   int `json:"all_count"`
	Unassigned int `json:"unassigned"`
}

// UpdateConversationRequest is the payload for updating conversation custom attributes.
type UpdateConversationRequest struct {
	CustomAttributes map[string]any `json:"custom_attributes"`
}

// TogglePriorityRequest is the payload for setting conversation priority.
type TogglePriorityRequest struct {
	Priority string `json:"priority"`
}

// ToggleStatusRequest is the payload for toggling conversation status.
type ToggleStatusRequest struct {
	Status string `json:"status"`
}

// AssignConversationRequest is the payload for assigning a conversation.
type AssignConversationRequest struct {
	AssigneeID *int `json:"assignee_id"`
	TeamID     *int `json:"team_id,omitempty"`
}

// ConversationLabelsRequest is the payload for updating conversation labels.
type ConversationLabelsRequest struct {
	Labels []string `json:"labels"`
}

// ---------------------------------------------------------------------------
// Message
// ---------------------------------------------------------------------------

// Message represents a Chatwoot message.
type Message struct {
	ID          int     `json:"id"`
	Content     *string `json:"content"`
	MessageType int     `json:"message_type"`
	ContentType string  `json:"content_type"`
	Private     bool    `json:"private"`
	Sender      *Sender `json:"sender"`
	CreatedAt   int64   `json:"created_at"`
}

// Sender represents who sent a message.
type Sender struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// MessageListResponse is the response from list messages.
type MessageListResponse struct {
	Payload []Message `json:"payload"`
}

// SendMessageRequest is the payload for sending a message.
type SendMessageRequest struct {
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
	Private     bool   `json:"private,omitempty"`
}

// DeleteMessageResponse is the response after deleting a message.
type DeleteMessageResponse struct {
	Success bool `json:"success"`
}

// ---------------------------------------------------------------------------
// Contact
// ---------------------------------------------------------------------------

// Contact represents a Chatwoot contact.
type Contact struct {
	ID               int            `json:"id"`
	Name             string         `json:"name"`
	Email            *string        `json:"email"`
	PhoneNumber      *string        `json:"phone_number"`
	Identifier       *string        `json:"identifier"`
	CreatedAt        time.Time      `json:"created_at"`
	LastActivityAt   *time.Time     `json:"last_activity_at"`
	CustomAttributes map[string]any `json:"custom_attributes"`
}

// ContactListResponse is the response from list contacts.
type ContactListResponse struct {
	Payload []Contact      `json:"payload"`
	Meta    PaginationMeta `json:"meta"`
}

// ContactSearchResponse is the response from search contacts.
type ContactSearchResponse struct {
	Payload []Contact `json:"payload"`
}

// CreateContactRequest is the payload for creating a contact.
type CreateContactRequest struct {
	InboxID          int            `json:"inbox_id,omitempty"`
	Name             string         `json:"name"`
	Email            string         `json:"email,omitempty"`
	PhoneNumber      string         `json:"phone_number,omitempty"`
	Identifier       string         `json:"identifier,omitempty"`
	CustomAttributes map[string]any `json:"custom_attributes,omitempty"`
}

// UpdateContactRequest is the payload for updating a contact.
type UpdateContactRequest struct {
	Name             *string        `json:"name,omitempty"`
	Email            *string        `json:"email,omitempty"`
	PhoneNumber      *string        `json:"phone_number,omitempty"`
	Identifier       *string        `json:"identifier,omitempty"`
	CustomAttributes map[string]any `json:"custom_attributes,omitempty"`
}

// ContactFilterRequest is the payload for filtering contacts.
type ContactFilterRequest struct {
	Payload []ContactFilterPayload `json:"payload"`
	Page    *int                   `json:"page,omitempty"`
}

// ContactFilterPayload describes a single filter condition for contacts.
type ContactFilterPayload struct {
	AttributeKey   string   `json:"attribute_key"`
	FilterOperator string   `json:"filter_operator"`
	Values         []string `json:"values"`
	QueryOperator  *string  `json:"query_operator,omitempty"`
}

// ContactConversationsResponse is the response listing conversations for a contact.
type ContactConversationsResponse struct {
	Payload []Conversation `json:"payload"`
}

// MergeContactsRequest is the payload for merging two contacts.
type MergeContactsRequest struct {
	BaseContactID   int `json:"base_contact_id"`
	MergeeContactID int `json:"mergee_contact_id"`
}

// ---------------------------------------------------------------------------
// Inbox
// ---------------------------------------------------------------------------

// Inbox represents a Chatwoot inbox.
type Inbox struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ChannelType string `json:"channel_type"`
	AvatarURL   string `json:"avatar_url"`
	Enabled     bool   `json:"enabled"`
}

// InboxListResponse is the response from list inboxes.
type InboxListResponse struct {
	Payload []Inbox `json:"payload"`
}

// ---------------------------------------------------------------------------
// Agent
// ---------------------------------------------------------------------------

// Agent represents a Chatwoot agent.
type Agent struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	AvailabilityStatus string `json:"availability_status"`
}

// ---------------------------------------------------------------------------
// Label
// ---------------------------------------------------------------------------

// Label represents a Chatwoot label.
type Label struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Color         string `json:"color"`
	ShowOnSidebar bool   `json:"show_on_sidebar"`
}

// ---------------------------------------------------------------------------
// Team
// ---------------------------------------------------------------------------

// Team represents a Chatwoot team.
type Team struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ---------------------------------------------------------------------------
// Canned Responses
// ---------------------------------------------------------------------------

// CannedResponse represents a saved canned response.
type CannedResponse struct {
	ID        int    `json:"id"`
	ShortCode string `json:"short_code"`
	Content   string `json:"content"`
}

// CreateCannedResponseRequest is the payload for creating a canned response.
type CreateCannedResponseRequest struct {
	ShortCode string `json:"short_code"`
	Content   string `json:"content"`
}

// UpdateCannedResponseRequest is the payload for updating a canned response.
type UpdateCannedResponseRequest struct {
	ShortCode *string `json:"short_code,omitempty"`
	Content   *string `json:"content,omitempty"`
}

// ---------------------------------------------------------------------------
// Custom Attributes
// ---------------------------------------------------------------------------

// CustomAttributeDefinition represents a custom attribute definition.
// AttributeModel: 0 = conversation_attribute, 1 = contact_attribute.
type CustomAttributeDefinition struct {
	ID                   int    `json:"id"`
	AttributeDisplayName string `json:"attribute_display_name"`
	AttributeDisplayType string `json:"attribute_display_type"`
	AttributeDescription string `json:"attribute_description"`
	AttributeKey         string `json:"attribute_key"`
	AttributeModel       int    `json:"attribute_model"`
}

// CreateCustomAttributeRequest is the payload for creating a custom attribute definition.
type CreateCustomAttributeRequest struct {
	AttributeDisplayName string `json:"attribute_display_name"`
	AttributeDisplayType string `json:"attribute_display_type"`
	AttributeDescription string `json:"attribute_description"`
	AttributeKey         string `json:"attribute_key"`
	AttributeModel       int    `json:"attribute_model"`
}

// ---------------------------------------------------------------------------
// Custom Filters
// ---------------------------------------------------------------------------

// CustomFilter represents a saved custom filter.
type CustomFilter struct {
	ID         int              `json:"id"`
	Name       string           `json:"name"`
	FilterType string           `json:"filter_type"`
	Query      json.RawMessage  `json:"query"`
}

// CreateCustomFilterRequest is the payload for creating a custom filter.
type CreateCustomFilterRequest struct {
	Name       string          `json:"name"`
	FilterType string          `json:"filter_type"`
	Query      json.RawMessage `json:"query"`
}

// ---------------------------------------------------------------------------
// Automation Rules
// ---------------------------------------------------------------------------

// AutomationRule represents an automation rule.
type AutomationRule struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	EventName   string           `json:"event_name"`
	Conditions  json.RawMessage  `json:"conditions"`
	Actions     json.RawMessage  `json:"actions"`
	Active      bool             `json:"active"`
	CreatedAt   time.Time        `json:"created_at"`
}

// CreateAutomationRuleRequest is the payload for creating an automation rule.
type CreateAutomationRuleRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	EventName   string          `json:"event_name"`
	Conditions  json.RawMessage `json:"conditions"`
	Actions     json.RawMessage `json:"actions"`
	Active      *bool           `json:"active,omitempty"`
}

// ---------------------------------------------------------------------------
// Webhooks
// ---------------------------------------------------------------------------

// Webhook represents a configured webhook.
type Webhook struct {
	ID            int      `json:"id"`
	URL           string   `json:"url"`
	Subscriptions []string `json:"subscriptions"`
}

// CreateWebhookRequest is the payload for creating a webhook.
type CreateWebhookRequest struct {
	URL           string   `json:"url"`
	Subscriptions []string `json:"subscriptions,omitempty"`
}

// UpdateWebhookRequest is the payload for updating a webhook.
type UpdateWebhookRequest struct {
	URL           *string  `json:"url,omitempty"`
	Subscriptions []string `json:"subscriptions,omitempty"`
}

// ---------------------------------------------------------------------------
// Reports
// ---------------------------------------------------------------------------

// ReportSummary contains aggregate metrics for a report.
type ReportSummary struct {
	AvgFirstResponseTime  float64 `json:"avg_first_response_time"`
	AvgResolutionTime     float64 `json:"avg_resolution_time"`
	ConversationsCount    int     `json:"conversations_count"`
	IncomingMessagesCount int     `json:"incoming_messages_count"`
	OutgoingMessagesCount int     `json:"outgoing_messages_count"`
	ResolutionsCount      int     `json:"resolutions_count"`
}

// AgentSummary contains report metrics for a single agent.
type AgentSummary struct {
	ID                    int     `json:"id"`
	Name                  string  `json:"name"`
	Email                 string  `json:"email"`
	AvgFirstResponseTime  float64 `json:"avg_first_response_time"`
	AvgResolutionTime     float64 `json:"avg_resolution_time"`
	ConversationsCount    int     `json:"conversations_count"`
	IncomingMessagesCount int     `json:"incoming_messages_count"`
	OutgoingMessagesCount int     `json:"outgoing_messages_count"`
	ResolutionsCount      int     `json:"resolutions_count"`
}

// TeamSummary contains report metrics for a single team.
type TeamSummary struct {
	ID                    int     `json:"id"`
	Name                  string  `json:"name"`
	AvgFirstResponseTime  float64 `json:"avg_first_response_time"`
	AvgResolutionTime     float64 `json:"avg_resolution_time"`
	ConversationsCount    int     `json:"conversations_count"`
	IncomingMessagesCount int     `json:"incoming_messages_count"`
	OutgoingMessagesCount int     `json:"outgoing_messages_count"`
	ResolutionsCount      int     `json:"resolutions_count"`
}

// ---------------------------------------------------------------------------
// Profile
// ---------------------------------------------------------------------------

// Profile represents the authenticated user's profile.
type Profile struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	AccountID int    `json:"account_id"`
	AvatarURL string `json:"avatar_url"`
}

// ---------------------------------------------------------------------------
// Team Members
// ---------------------------------------------------------------------------

// TeamMember represents a member of a team (same structure as Agent).
type TeamMember struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	AvailabilityStatus string `json:"availability_status"`
}

// AddTeamMemberRequest is the payload for adding members to a team.
type AddTeamMemberRequest struct {
	UserIDs []int `json:"user_ids"`
}

// ---------------------------------------------------------------------------
// Inbox Members
// ---------------------------------------------------------------------------

// InboxMember represents a member of an inbox (same structure as Agent).
type InboxMember struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	AvailabilityStatus string `json:"availability_status"`
}

// AddInboxMemberRequest is the payload for adding members to an inbox.
type AddInboxMemberRequest struct {
	UserIDs []int `json:"user_ids"`
}

// ---------------------------------------------------------------------------
// Help Center - Portals
// ---------------------------------------------------------------------------

// Portal represents a help center portal.
type Portal struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ---------------------------------------------------------------------------
// Help Center - Articles
// ---------------------------------------------------------------------------

// Article represents a help center article.
type Article struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	Content    *string `json:"content"`
	Status     string  `json:"status"`
	PortalID   int     `json:"portal_id"`
	CategoryID *int    `json:"category_id"`
	AuthorID   *int    `json:"author_id"`
}

// ---------------------------------------------------------------------------
// Help Center - Categories
// ---------------------------------------------------------------------------

// Category represents a help center category.
type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	PortalID int    `json:"portal_id"`
}

// ---------------------------------------------------------------------------
// Notifications
// ---------------------------------------------------------------------------

// Notification represents a Chatwoot notification.
type Notification struct {
	ID               int        `json:"id"`
	NotificationType string     `json:"notification_type"`
	PrimaryActorType string     `json:"primary_actor_type"`
	PrimaryActorID   int        `json:"primary_actor_id"`
	ReadAt           *time.Time `json:"read_at"`
	CreatedAt        time.Time  `json:"created_at"`
}
