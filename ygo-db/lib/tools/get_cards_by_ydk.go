package tools

import (
	"context"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/cdb"
	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/ydk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetCardsByYDKInput represents the input parameters for the get_cards_by_ydk tool.
type GetCardsByYDKInput struct {
	// YDK is the content of a YDK file (Yugioh Deck file)
	YDK string `json:"ydk" jsonschema:"the content of a YDK file (Yugioh Deck file)"`
}

// GetCardsByYDKOutput represents the output of the get_cards_by_ydk tool.
type GetCardsByYDKOutput struct {
	Main  []*cdb.CardInfoForAI `json:"main"`
	Extra []*cdb.CardInfoForAI `json:"extra"`
	Side  []*cdb.CardInfoForAI `json:"side"`
}

// GetCardsByYDK retrieves cards from a YDK file content.
func GetCardsByYDK(db *cdb.DB, ydkContent string) ([]*cdb.CardInfoForAI, []*cdb.CardInfoForAI, []*cdb.CardInfoForAI, error) {
	// Parse YDK file to get main, extra, side decks
	mainDeck, extraDeck, sideDeck := ydk.ParseYDKFile(ydkContent)

	// Convert map keys to []uint64 for GetCardsByIDs
	mainIDs := make([]uint64, 0, len(mainDeck))
	for code := range mainDeck {
		mainIDs = append(mainIDs, uint64(code))
	}

	extraIDs := make([]uint64, 0, len(extraDeck))
	for code := range extraDeck {
		extraIDs = append(extraIDs, uint64(code))
	}

	sideIDs := make([]uint64, 0, len(sideDeck))
	for code := range sideDeck {
		sideIDs = append(sideIDs, uint64(code))
	}

	// Call GetCardsByIDs 3 times for main, extra, side
	mainCards, err := db.GetCardsByIDs(mainIDs)
	if err != nil {
		return nil, nil, nil, err
	}

	extraCards, err := db.GetCardsByIDs(extraIDs)
	if err != nil {
		return nil, nil, nil, err
	}

	sideCards, err := db.GetCardsByIDs(sideIDs)
	if err != nil {
		return nil, nil, nil, err
	}

	// Convert to CardInfoForAI, using Count field instead of duplicating cards
	mainAI := make([]*cdb.CardInfoForAI, 0, len(mainDeck))
	for code, count := range mainDeck {
		card, ok := mainCards[uint64(code)]
		if ok {
			ai := card.ToCardInfoForAI()
			ai.Count = count
			mainAI = append(mainAI, ai)
		}
	}

	extraAI := make([]*cdb.CardInfoForAI, 0, len(extraDeck))
	for code, count := range extraDeck {
		card, ok := extraCards[uint64(code)]
		if ok {
			ai := card.ToCardInfoForAI()
			ai.Count = count
			extraAI = append(extraAI, ai)
		}
	}

	sideAI := make([]*cdb.CardInfoForAI, 0, len(sideDeck))
	for code, count := range sideDeck {
		card, ok := sideCards[uint64(code)]
		if ok {
			ai := card.ToCardInfoForAI()
			ai.Count = count
			sideAI = append(sideAI, ai)
		}
	}

	return mainAI, extraAI, sideAI, nil
}

// GetCardsByYDKHandler creates a tool handler for getting cards by YDK file.
func GetCardsByYDKHandler(db *cdb.DB) func(ctx context.Context, req *mcp.CallToolRequest, args GetCardsByYDKInput) (*mcp.CallToolResult, *GetCardsByYDKOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args GetCardsByYDKInput) (*mcp.CallToolResult, *GetCardsByYDKOutput, error) {
		if args.YDK == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "ydk content is required"},
				},
			}, nil, nil
		}

		main, extra, side, err := GetCardsByYDK(db, args.YDK)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: err.Error()},
				},
			}, nil, nil
		}

		result := &GetCardsByYDKOutput{
			Main:  main,
			Extra: extra,
			Side:  side,
		}

		return nil, result, nil
	}
}
