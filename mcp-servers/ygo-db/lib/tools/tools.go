package tools

import (
	"github.com/elmhuangyu/yu-gi-oh-ai-tools/mcp-servers/ygo-db/lib/cdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolNames defines the names of all available tools.
const (
	ToolNameGetCardByID          = "get_card_by_id"
	ToolNameGetCardsByName       = "get_cards_by_name"
	ToolNameGetCardsByArchetypes = "get_cards_by_archetypes"
)

// Tools returns a slice of all MCP tools.
func Tools(server *mcp.Server, db *cdb.DB) {
	// Register get_card_by_id tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        ToolNameGetCardByID,
		Description: "Retrieve a Yu-Gi-Oh! card by its unique passcode (ID). Takes a card ID as input and returns the card's information if found.",
	}, GetCardByIDHandler(db))

	// Register get_cards_by_name tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        ToolNameGetCardsByName,
		Description: "Search for Yu-Gi-Oh! cards by name. Takes a card name and page number as input. Returns exact match (if on page 0), partial matches, total count, and current page.",
	}, GetCardsByNameHandler(db))

	// Register get_cards_by_archetypes tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        ToolNameGetCardsByArchetypes,
		Description: "Search for Yu-Gi-Oh! cards by archetype/set names. Takes a list of archetype names (1-4) and page number as input. Returns matching cards, total count, and current page.",
	}, GetCardsByArchetypesHandler(db))

}
