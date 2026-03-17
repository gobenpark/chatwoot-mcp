# chatwoot-mcp

A Model Context Protocol (MCP) server for [Chatwoot](https://www.chatwoot.com), built in Go. Connect Claude (or any MCP client) directly to your Chatwoot instance to manage conversations, contacts, reports, and more through natural language.

## Features

| Category | Tool | Description |
|----------|------|-------------|
| **Conversations** | `list_conversations` | List conversations with status filter |
| | `get_conversation` | Get detailed conversation info |
| | `filter_conversations` | Filter with advanced criteria |
| | `get_conversation_counts` | Get counts by status |
| | `get_messages` | Get all messages in a conversation |
| | `send_message` | Send outgoing, incoming, or private note |
| | `delete_message` | Delete a message |
| | `toggle_conversation_status` | Change status (open, resolved, pending, snoozed) |
| | `toggle_conversation_priority` | Set priority (urgent, high, medium, low) |
| | `assign_conversation` | Assign to agent and/or team |
| | `update_conversation_labels` | Update conversation labels |
| **Contacts** | `list_contacts` | List contacts with pagination |
| | `get_contact` | Get contact details |
| | `search_contacts` | Search by name, email, phone |
| | `create_contact` | Create a new contact |
| | `update_contact` | Update contact info |
| | `delete_contact` | Delete a contact |
| | `filter_contacts` | Filter with advanced criteria |
| | `get_contact_conversations` | List conversations for a contact |
| | `merge_contacts` | Merge two contacts |
| | `update_contact_labels` | Update contact labels |
| **Account** | `get_account` | Get account details and settings |
| | `update_account` | Update account settings |
| | `list_inboxes` | List all inboxes (channels) |
| | `list_agents` | List all agents |
| | `list_labels` | List all labels |
| | `list_teams` | List all teams |
| | `get_team` | Get team details |
| | `list_team_members` | List agents in a team |
| | `list_inbox_members` | List agents in an inbox |
| | `get_profile` | Get authenticated user profile |
| **Reports** | `get_reports_summary` | Account-level metrics with period comparison |
| | `get_agent_summary` | Per-agent performance metrics |
| | `get_team_summary` | Per-team performance metrics |
| | `get_inbox_summary` | Per-inbox performance metrics |
| | `get_channel_summary` | Per-channel type metrics |
| **Automation** | `list_canned_responses` | List saved reply templates |
| | `create_canned_response` | Create a canned response |
| | `delete_canned_response` | Delete a canned response |
| | `list_custom_attributes` | List custom attribute definitions |
| | `list_custom_filters` | List saved custom filters |
| | `list_automation_rules` | List automation rules |
| | `list_webhooks` | List webhooks |
| | `create_webhook` | Create a webhook |
| | `delete_webhook` | Delete a webhook |
| **Agent Bots** | `list_agent_bots` | List all agent bots |
| | `get_agent_bot` | Get agent bot details |
| | `create_agent_bot` | Create a new agent bot |
| | `update_agent_bot` | Update an agent bot |
| | `delete_agent_bot` | Delete an agent bot |
| **Integrations** | `list_integration_apps` | List available integrations |
| | `create_integration_hook` | Create an integration hook |
| | `update_integration_hook` | Update an integration hook |
| | `delete_integration_hook` | Delete an integration hook |
| **Help Center** | `list_portals` | List help center portals |
| | `update_portal` | Update a portal |
| | `list_articles` | List articles in a portal |
| | `create_article` | Create an article |
| | `update_article` | Update an article |
| | `delete_article` | Delete an article |
| | `list_categories` | List categories in a portal |
| | `create_category` | Create a category |
| | `update_category` | Update a category |
| | `delete_category` | Delete a category |
| **Audit** | `list_audit_logs` | List audit log entries |

## Installation

### npx (Recommended)

No prerequisites required. The binary is downloaded automatically on first run:

```bash
npx chatwoot-mcp
```

### Pre-built binaries

Download from [GitHub Releases](https://github.com/gobenpark/chatwoot-mcp/releases):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/gobenpark/chatwoot-mcp/releases/latest/download/chatwoot-mcp_darwin_arm64.tar.gz | tar xz
mv chatwoot-mcp /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/gobenpark/chatwoot-mcp/releases/latest/download/chatwoot-mcp_darwin_amd64.tar.gz | tar xz
mv chatwoot-mcp /usr/local/bin/

# Linux (amd64)
curl -L https://github.com/gobenpark/chatwoot-mcp/releases/latest/download/chatwoot-mcp_linux_amd64.tar.gz | tar xz
mv chatwoot-mcp /usr/local/bin/
```

### From source

```bash
go install github.com/gobenpark/chatwoot-mcp@latest
```

## Configuration

### 1. Get a Chatwoot API Token

Go to **Chatwoot** > **Profile Settings** > **Access Token** and copy your API access token.

### 2. Configure Claude Code

Using the CLI:

```bash
# Add to current project
claude mcp add chatwoot \
  -e CHATWOOT_URL=https://your-chatwoot.example.com \
  -e CHATWOOT_API_TOKEN=your_api_access_token \
  -e CHATWOOT_ACCOUNT_ID=1 \
  -- npx -y chatwoot-mcp

# Add globally (available in all projects)
claude mcp add chatwoot -s user \
  -e CHATWOOT_URL=https://your-chatwoot.example.com \
  -e CHATWOOT_API_TOKEN=your_api_access_token \
  -e CHATWOOT_ACCOUNT_ID=1 \
  -- npx -y chatwoot-mcp
```

Or manually add to your `.mcp.json` (project-level) or `~/.claude/settings.json` (global):

```json
{
  "mcpServers": {
    "chatwoot": {
      "command": "npx",
      "args": ["-y", "chatwoot-mcp"],
      "env": {
        "CHATWOOT_URL": "https://your-chatwoot.example.com",
        "CHATWOOT_API_TOKEN": "your_api_access_token",
        "CHATWOOT_ACCOUNT_ID": "1"
      }
    }
  }
}
```

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `CHATWOOT_URL` | Yes | Your Chatwoot instance URL |
| `CHATWOOT_API_TOKEN` | Yes | API access token from Profile Settings |
| `CHATWOOT_ACCOUNT_ID` | Yes | Your Chatwoot account ID |

## Usage Examples

Once configured, use natural language with Claude:

```
> List all open conversations
> Show me the messages in conversation #273
> Get a summary of agent performance this week
> Search for contacts with email containing "gmail.com"
> Send a private note to conversation #292
> Create a new contact named "John Doe" with email john@example.com
> What are the conversation counts by status?
> List all automation rules
> Show me the inbox performance report for this month
```

## Development

```bash
# Clone
git clone https://github.com/gobenpark/chatwoot-mcp.git
cd chatwoot-mcp

# Build
go build -o chatwoot-mcp .

# Test locally
export CHATWOOT_URL="https://your-chatwoot.example.com"
export CHATWOOT_API_TOKEN="your_token"
export CHATWOOT_ACCOUNT_ID="1"
./chatwoot-mcp
```

## Release

Releases are automated via GitHub Actions + [GoReleaser](https://goreleaser.com). To create a release:

```bash
git tag v0.3.0
git push origin v0.3.0
```

This builds binaries for Linux, macOS, and Windows (amd64/arm64) and publishes to npm automatically.

## License

MIT
