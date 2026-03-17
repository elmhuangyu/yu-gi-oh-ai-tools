package tools

import (
	"context"
	"errors"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/cdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetCardsByArchetypesInput represents the input parameters for the get_cards_by_archetypes tool.
type GetCardsByArchetypesInput struct {
	// Archetypes is the list of archetype names to search for (1-4 archetypes)
	Archetypes []string `json:"archetypes" jsonschema:"the list of archetype names to search for (1-4 archetypes)"`
	// Page is the page number (0-indexed)
	Page int `json:"page" jsonschema:"the page number (0-indexed)"`
}

// GetCardsByArchetypesOutput represents the output of the get_cards_by_archetypes tool.
type GetCardsByArchetypesOutput struct {
	Cards       []*cdb.CardInfoForAI `json:"cards"`
	Total       int                  `json:"total"`
	CurrentPage int                  `json:"currentPage"`
}

// GetCardsByArchetypes retrieves cards by archetype/set names using the provided database.
func GetCardsByArchetypes(db *cdb.DB, archetypes []string, page int) ([]*cdb.CardInfoForAI, int, error) {
	const pageSize = 30
	cards, total, err := db.FindCardsBySetName(archetypes, pageSize, page)
	if err != nil {
		return nil, 0, err
	}

	if len(cards) == 0 {
		return nil, 0, errors.New("no cards found for the specified archetypes")
	}

	cardsAI := make([]*cdb.CardInfoForAI, len(cards))
	for i, card := range cards {
		cardsAI[i] = card.ToCardInfoForAI()
	}

	return cardsAI, total, nil
}

// GetCardsByArchetypesHandler creates a tool handler for getting cards by archetypes.
func GetCardsByArchetypesHandler(db *cdb.DB) func(ctx context.Context, req *mcp.CallToolRequest, args GetCardsByArchetypesInput) (*mcp.CallToolResult, *GetCardsByArchetypesOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args GetCardsByArchetypesInput) (*mcp.CallToolResult, *GetCardsByArchetypesOutput, error) {
		if len(args.Archetypes) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "at least one archetype is required"},
				},
			}, nil, nil
		}

		if len(args.Archetypes) > 4 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "maximum 4 archetypes allowed"},
				},
			}, nil, nil
		}

		if args.Page < 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "page must be >= 0"},
				},
			}, nil, nil
		}

		cards, total, err := GetCardsByArchetypes(db, args.Archetypes, args.Page)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: err.Error()},
				},
			}, nil, nil
		}

		result := &GetCardsByArchetypesOutput{
			Cards:       cards,
			Total:       total,
			CurrentPage: args.Page,
		}

		return nil, result, nil
	}
}
