package commands

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/cdb"
)

// GetCardsByArchetypesCmd represents the get-cards-by-archetypes command
var GetCardsByArchetypesCmd = &cobra.Command{
	Use:   "get-cards-by-archetypes",
	Short: `Search for cards by archetype names (1-4 archetypes). You should use "a archetype" if archetype name contains more than one word.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		archetypes := args

		if len(archetypes) == 0 {
			log.Fatal("at least one archetype is required")
		}

		if len(archetypes) > 4 {
			log.Fatal("maximum 4 archetypes allowed")
		}

		if pageFlag < 0 {
			log.Fatal("page must be >= 0")
		}

		cardsForAI, cardsForHuman, total, err := getCardsByArchetypes(db, archetypes, pageFlag)
		if err != nil {
			log.Fatalf("failed to get cards: %v", err)
		}

		if len(cardsForAI) == 0 {
			fmt.Printf("No cards found for the specified archetypes: %v\n", archetypes)
			return
		}

		result := &getCardsByArchetypesOutput{
			Cards:       cardsForAI,
			Total:       total,
			CurrentPage: pageFlag,
		}

		// if we return all, no need for page information
		if len(cardsForAI) == total {
			result.CurrentPage = 0
			result.Total = 0
		} else if outputFormat == "csv" {
			fmt.Println("csv does not support pagination. Consider increase the page size.")
			return
		}

		output, err := formatArchetypesOutput(result, cardsForHuman, outputFormat)
		if err != nil {
			log.Fatalf("failed to format output: %v", err)
		}
		fmt.Println(string(output))
	},
}

type getCardsByArchetypesOutput struct {
	Cards       []*cdb.CardInfoForAI `json:"cards" yaml:"cards"`
	Total       int                  `json:"total,omitempty" yaml:"total,omitempty"`
	CurrentPage int                  `json:"currentPage,omitempty" yaml:"currentPage,omitempty"`
}

func getCardsByArchetypes(database *cdb.DB, archetypes []string, page int) ([]*cdb.CardInfoForAI, []*cdb.CardInfoForHuman, int, error) {
	// 1 archetype usually lower than 200. for example: HERO is ~150.
	const pageSize = 250
	cards, total, err := database.FindCardsBySetName(archetypes, pageSize, page)
	if err != nil {
		return nil, nil, 0, err
	}

	if len(cards) == 0 {
		return nil, nil, 0, nil
	}

	cardsAI := make([]*cdb.CardInfoForAI, len(cards))
	for i, card := range cards {
		cardsAI[i] = card.ToCardInfoForAI()
	}

	return cardsAI, cards, total, nil
}

// formatArchetypesOutput formats the getCardsByArchetypesOutput to the specified format.
func formatArchetypesOutput(result *getCardsByArchetypesOutput, cardsForHuman []*cdb.CardInfoForHuman, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.Marshal(result)
	case "yaml":
		return yaml.Marshal(result)
	case "csv":
		return formatCardsForHumanToCSV(cardsForHuman)
	default:
		return nil, fmt.Errorf("unsupported output format: %s (supported: json, yaml, csv)", format)
	}
}

func init() {
	GetCardsByArchetypesCmd.Flags().StringVar(&outputFormat, "format", "yaml", "Output format: json, yaml, csv")
	GetCardsByArchetypesCmd.Flags().IntVar(&pageFlag, "page", 0, "Page number (0-indexed)")
}
