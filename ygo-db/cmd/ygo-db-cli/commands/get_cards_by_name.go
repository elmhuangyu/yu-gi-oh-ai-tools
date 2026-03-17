package commands

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/cdb"
)

var (
	pageFlag        int
	defaultPageSize = 30
)

// GetCardsByNameCmd represents the get-cards-by-name command
var GetCardsByNameCmd = &cobra.Command{
	Use:   "get-cards-by-name",
	Short: "Search for cards by name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if name == "" {
			log.Fatal("name is required")
		}

		if pageFlag < 0 {
			log.Fatal("page must be >= 0")
		}

		exactAI, maybeAI, exactHuman, maybeHuman, total, err := getCardsByName(db, name, pageFlag)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Card with name %q not found", name)
				os.Exit(0)
			}
			log.Fatalf("failed to get cards: %v", err)
		}

		// Build cards for output based on format
		var cardsForAI []*cdb.CardInfoForAI
		if exactAI != nil {
			cardsForAI = append(cardsForAI, exactAI)
		}
		cardsForAI = append(cardsForAI, maybeAI...)

		result := &getCardsByNameOutput{
			Exact:       exactAI,
			Maybe:       maybeAI,
			Total:       total,
			CurrentPage: pageFlag,
		}

		// Check if we have all results (for pagination handling)
		// If total > defaultPageSize, there are more results than we can show
		if total > defaultPageSize && outputFormat == "csv" {
			log.Fatalf("csv does not support pagination. Found %d results, but limit is %d. Consider using a more specific search term or increase page size.", total, defaultPageSize)
		}

		// if we return all, no need for page information
		if len(cardsForAI) == total {
			result.CurrentPage = 0
			result.Total = 0
		} else if outputFormat == "csv" {
			log.Fatalf("csv does not support pagination. Found %d results, but only showing %d. Consider using a more specific search term.", total, len(cardsForAI))
		}

		output, err := formatOutput(result, exactHuman, maybeHuman, outputFormat)
		if err != nil {
			log.Fatalf("failed to format output: %v", err)
		}
		fmt.Println(string(output))
	},
}

type getCardsByNameOutput struct {
	Exact       *cdb.CardInfoForAI   `json:"exact" yaml:"exact"`
	Maybe       []*cdb.CardInfoForAI `json:"maybe" yaml:"maybe"`
	Total       int                  `json:"total" yaml:"total"`
	CurrentPage int                  `json:"currentPage" yaml:"currentPage"`
}

func getCardsByName(database *cdb.DB, name string, page int) (*cdb.CardInfoForAI, []*cdb.CardInfoForAI, *cdb.CardInfoForHuman, []*cdb.CardInfoForHuman, int, error) {
	exact, maybe, total, err := database.FindCardByName(name, page)
	if err != nil {
		return nil, nil, nil, nil, 0, err
	}

	var exactAI *cdb.CardInfoForAI
	var exactHuman *cdb.CardInfoForHuman
	if exact != nil {
		exactAI = exact.ToCardInfoForAI()
		exactHuman = exact
	}

	maybeAI := make([]*cdb.CardInfoForAI, len(maybe))
	maybeHuman := make([]*cdb.CardInfoForHuman, len(maybe))
	for i, card := range maybe {
		maybeAI[i] = card.ToCardInfoForAI()
		maybeHuman[i] = card
	}

	return exactAI, maybeAI, exactHuman, maybeHuman, total, nil
}

// formatOutput formats the getCardsByNameOutput to the specified format.
func formatOutput(result *getCardsByNameOutput, exactHuman *cdb.CardInfoForHuman, maybeHuman []*cdb.CardInfoForHuman, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.Marshal(result)
	case "yaml":
		return yaml.Marshal(result)
	case "csv":
		// Combine exact and maybe into a single slice for CSV output
		var allCards []*cdb.CardInfoForHuman
		if exactHuman != nil {
			allCards = append(allCards, exactHuman)
		}
		allCards = append(allCards, maybeHuman...)
		return formatCardsForHumanToCSV(allCards)
	default:
		return nil, fmt.Errorf("unsupported output format: %s (supported: json, yaml, csv)", format)
	}
}

func init() {
	GetCardsByNameCmd.Flags().StringVar(&outputFormat, "format", "yaml", "Output format: json, yaml, csv")
	GetCardsByNameCmd.Flags().IntVar(&pageFlag, "page", 0, "Page number (0-indexed)")
}
