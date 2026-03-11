package tools

import (
	"github.com/elmhuangyu/yu-gi-oh-ai-tools/mcp-servers/ygo-db/lib/cdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolNames defines the names of all available tools.
const ToolNameGetCardByID = "get_card_by_id"

// Tools returns a slice of all MCP tools.
func Tools(server *mcp.Server, db *cdb.DB) {
	// Register get_card_by_id tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        ToolNameGetCardByID,
		Description: "Retrieve a Yu-Gi-Oh! card by its unique passcode (ID). Takes a card ID as input and returns the card's information if found.",
	}, GetCardByIDHandler(db))

}
