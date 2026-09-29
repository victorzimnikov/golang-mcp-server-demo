package mcpserver

import (
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

func getSourceFromRequest(request *mcp.CallToolRequest) domain.Source {
	source := domain.SourceHuman

	clientInfo := request.ClientInfo()

	if clientInfo != nil {
		clientName := strings.ToLower(clientInfo.Name)

		if strings.Contains(clientName, "codex") {
			source = domain.SourceCodex
		} else if strings.Contains(clientName, "claude") {
			source = domain.SourceClaude
		}
	}

	return source
}
