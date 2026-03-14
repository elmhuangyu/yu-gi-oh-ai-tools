package tools

import (
	"context"
	"database/sql"
	"errors"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/cdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetCardsByNameInput represents the input parameters for the get_cards_by_name tool.
type GetCardsByNameInput struct {
	// Name is the card name to search for
	Name string `json:"name" jsonschema:"the card name to search for"`
	// Page is the page number (0-indexed)
	Page int `json:"page" jsonschema:"the page number (0-indexed)"`
}

// GetCardsByNameOutput represents the output of the get_cards_by_name tool.
type GetCardsByNameOutput struct {
	Exact       *cdb.CardInfoForAI   `json:"exact"`
	Maybe       []*cdb.CardInfoForAI `json:"maybe"`
	Total       int                  `json:"total"`
	CurrentPage int                  `json:"currentPage"`
}

// GetCardsByName retrieves cards by name using the provided database.
func GetCardsByName(db *cdb.DB, name string, page int) (*cdb.CardInfoForAI, []*cdb.CardInfoForAI, int, error) {
	exact, maybe, total, err := db.FindCardByName(name, page)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, 0, errors.New("card not found")
		}
		return nil, nil, 0, err
	}

	var exactAI *cdb.CardInfoForAI
	if exact != nil {
		exactAI = exact.ToCardInfoForAI()
	}

	maybeAI := make([]*cdb.CardInfoForAI, len(maybe))
	for i, card := range maybe {
		maybeAI[i] = card.ToCardInfoForAI()
	}

	return exactAI, maybeAI, total, nil
}

// GetCardsByNameHandler creates a tool handler for getting cards by name.
func GetCardsByNameHandler(db *cdb.DB) func(ctx context.Context, req *mcp.CallToolRequest, args GetCardsByNameInput) (*mcp.CallToolResult, *GetCardsByNameOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args GetCardsByNameInput) (*mcp.CallToolResult, *GetCardsByNameOutput, error) {
		if args.Name == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "name is required"},
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

		exact, maybe, total, err := GetCardsByName(db, args.Name, args.Page)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: err.Error()},
				},
			}, nil, nil
		}

		result := &GetCardsByNameOutput{
			Exact:       exact,
			Maybe:       maybe,
			Total:       total,
			CurrentPage: args.Page,
		}

		return nil, result, nil
	}
}
