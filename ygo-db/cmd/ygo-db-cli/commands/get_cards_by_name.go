package commands

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/cdb"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

var (
	pageFlag int
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

		exact, maybe, total, err := getCardsByName(db, name, pageFlag)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Card with name %q not found", name)
				os.Exit(0)
			}
			log.Fatalf("failed to get cards: %v", err)
		}

		result := &getCardsByNameOutput{
			Exact:       exact,
			Maybe:       maybe,
			Total:       total,
			CurrentPage: pageFlag,
		}

		output, err := formatOutput(result, outputFormat)
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

func getCardsByName(database *cdb.DB, name string, page int) (*cdb.CardInfoForAI, []*cdb.CardInfoForAI, int, error) {
	exact, maybe, total, err := database.FindCardByName(name, page)
	if err != nil {
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

// formatOutput formats the getCardsByNameOutput to the specified format.
func formatOutput(result *getCardsByNameOutput, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.Marshal(result)
	case "yaml":
		return yaml.Marshal(result)
	default:
		return nil, fmt.Errorf("unsupported output format: %s (supported: json, yaml)", format)
	}
}

func init() {
	GetCardsByNameCmd.Flags().StringVar(&outputFormat, "format", "yaml", "Output format: json, yaml")
	GetCardsByNameCmd.Flags().IntVar(&pageFlag, "page", 0, "Page number (0-indexed)")
}
