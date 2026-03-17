package commands

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/cdb"
	"github.com/moznion/go-optional"
)

// formatCardsForHumanToCSV converts a slice of CardInfoForHuman to CSV format.
func formatCardsForHumanToCSV(cards []*cdb.CardInfoForHuman) ([]byte, error) {
	if len(cards) == 0 {
		return []byte{}, nil
	}

	headers, rows := cdb.CardInfoForHumanToCSV(cards)
	if headers == nil {
		return []byte{}, nil
	}

	return buildCSV(headers, rows)
}

// buildCSV builds CSV output from headers and rows.
func buildCSV(headers []string, rows [][]string) ([]byte, error) {
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	err := writer.Write(headers)
	if err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	err = writer.WriteAll(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to write CSV row: %w", err)
	}

	writer.Flush()
	return []byte(buf.String()), nil
}

// cardInfoAItoHuman converts CardInfoForAI to CardInfoForHuman for CSV output.
// This is a package-level helper function for use by all command files.
func cardInfoAItoHuman(ai *cdb.CardInfoForAI) *cdb.CardInfoForHuman {
	human := &cdb.CardInfoForHuman{
		Name: ai.Name,
		Desc: ai.Desc,
	}

	if ai.Atk != nil {
		human.Atk = some(int(*ai.Atk))
	}
	if ai.Def != nil {
		human.Def = some(int(*ai.Def))
	}
	if ai.Level != nil {
		human.Level = some(int(*ai.Level))
	}
	if ai.Race != nil {
		human.Race = some(*ai.Race)
	}
	if ai.Attribute != nil {
		human.Attribute = some(*ai.Attribute)
	}
	if ai.Type != "" {
		human.Type = strings.Split(ai.Type, "|")
	}
	if ai.Archetypes != "" {
		human.SetNames = strings.Split(ai.Archetypes, "|")
	}

	return human
}

// some is a helper to create an optional value.
func some[T any](v T) optional.Option[T] {
	return optional.Some(v)
}
