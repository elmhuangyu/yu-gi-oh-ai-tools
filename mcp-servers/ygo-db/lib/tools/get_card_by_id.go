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

// GetCardByIDOutput represents the output of the get_card_by_id tool.
type GetCardByIDOutput struct {
	Card  *cdb.CardInfoForAI `json:"card,omitempty"`
	Found bool               `json:"found"`
}

// GetCardByIDHandler creates a tool handler for getting a card by ID.
func GetCardByIDHandler(db *cdb.DB) func(ctx context.Context, req *mcp.CallToolRequest, args GetCardByIDInput) (*mcp.CallToolResult, *GetCardByIDOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args GetCardByIDInput) (*mcp.CallToolResult, *GetCardByIDOutput, error) {
		if args.ID == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "id is required"},
				},
			}, nil, nil
		}

		card, err := db.GetCardByID(args.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, &GetCardByIDOutput{Found: false}, nil
			}
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: err.Error()},
				},
			}, nil, nil
		}

		return nil, &GetCardByIDOutput{Card: card.ToCardInfoForAI(), Found: true}, nil
	}
}
