package tools

import (
	"context"
	"database/sql"
	"errors"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/mcp-servers/ygo-db/lib/cdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetCardByIDInput represents the input parameters for the get_card_by_id tool.
type GetCardByIDInput struct {
	// ID is the Yu-Gi-Oh! card's unique passcode
	ID uint64 `json:"id" jsonschema:"the card's unique passcode"`
}

// GetCardByID retrieves a card by ID using the provided database.
func GetCardByID(db *cdb.DB, id uint64) (*cdb.CardInfoForAI, error) {
	card, err := db.GetCardByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("card not found")
		}
		return nil, err
	}
	return card.ToCardInfoForAI(), nil
}

// GetCardByIDHandler creates a tool handler for getting a card by ID.
func GetCardByIDHandler(db *cdb.DB) func(ctx context.Context, req *mcp.CallToolRequest, args GetCardByIDInput) (*mcp.CallToolResult, *cdb.CardInfoForAI, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args GetCardByIDInput) (*mcp.CallToolResult, *cdb.CardInfoForAI, error) {
		if args.ID == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "id is required"},
				},
			}, nil, nil
		}

		card, err := GetCardByID(db, args.ID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: err.Error()},
				},
			}, nil, nil
		}

		return nil, card, nil
	}
}
