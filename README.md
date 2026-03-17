# chatwoot-mcp

MCP (Model Context Protocol) server for [Chatwoot](https://www.chatwoot.com/) — connect Claude to your Chatwoot instance.

## Features

### Conversations
- `list_conversations` — List conversations with status filter (open, resolved, pending, snoozed)
- `get_conversation` — Get detailed conversation info
- `get_messages` — Get all messages in a conversation
- `send_message` — Send a message (outgoing, incoming, or private note)
- `toggle_conversation_status` — Change status (open, resolved, pending, snoozed)
- `assign_conversation` — Assign to agent and/or team
- `update_conversation_labels` — Update conversation labels

### Contacts
- `list_contacts` — List contacts with pagination
- `get_contact` — Get contact details
- `search_contacts` — Search by name, email, phone, or identifier
- `create_contact` — Create a new contact

### Account
- `list_inboxes` — List all inboxes
- `list_agents` — List all agents
- `list_labels` — List all labels
- `list_teams` — List all teams

## Installation

### From npm

```bash
npx chatwoot-mcp
```

### From source

```bash
go install github.com/gobenpark/chatwoot-mcp@latest
```

### Build from source

```bash
git clone https://github.com/gobenpark/chatwoot-mcp.git
cd chatwoot-mcp
go build -o chatwoot-mcp .
```

## Configuration

Set the following environment variables:

| Variable | Description | Required |
|----------|-------------|----------|
| `CHATWOOT_URL` | Your Chatwoot instance URL (e.g., `https://app.chatwoot.com`) | Yes |
| `CHATWOOT_API_TOKEN` | API access token from Profile Settings | Yes |
| `CHATWOOT_ACCOUNT_ID` | Your Chatwoot account ID | Yes |

### Claude Desktop / Claude Code

Add to your MCP configuration:

```json
{
  "mcpServers": {
    "chatwoot": {
      "command": "npx",
      "args": ["chatwoot-mcp"],
      "env": {
        "CHATWOOT_URL": "https://your-chatwoot.example.com",
        "CHATWOOT_API_TOKEN": "your_api_access_token",
        "CHATWOOT_ACCOUNT_ID": "1"
      }
    }
  }
}
```

## License

MIT
