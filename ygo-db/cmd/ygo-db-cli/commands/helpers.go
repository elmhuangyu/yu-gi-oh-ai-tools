package commands

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/cdb"
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
