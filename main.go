package main

import (
	"context"
	"log"
	"os"

	"github.com/gobenpark/chatwoot-mcp/chatwoot"
	"github.com/gobenpark/chatwoot-mcp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	log.SetOutput(os.Stderr)

	client, err := chatwoot.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize Chatwoot client: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "chatwoot-mcp",
		Version: version,
	}, nil)

	tools.RegisterConversationTools(server, client)
	tools.RegisterContactTools(server, client)
	tools.RegisterAccountTools(server, client)
	tools.RegisterReportTools(server, client)
	tools.RegisterAutomationTools(server, client)
	tools.RegisterAuditTools(server, client)
	tools.RegisterAgentBotTools(server, client)
	tools.RegisterIntegrationTools(server, client)
	tools.RegisterHelpCenterTools(server, client)

	log.Println("Starting chatwoot-mcp server (stdio)...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
