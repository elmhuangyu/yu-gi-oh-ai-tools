package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/cdb"
	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/ydk"
	"github.com/goccy/go-yaml"
	"github.com/moznion/go-optional"
	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
)

// getCardsByYDKOutput represents the output of the get-cards-by-ydk command.
type getCardsByYDKOutput struct {
	Main  []*cdb.CardInfoForAI `json:"main" yaml:"main"`
	Extra []*cdb.CardInfoForAI `json:"extra,omitempty" yaml:"extra,omitempty"`
	Side  []*cdb.CardInfoForAI `json:"side,omitempty" yaml:"side,omitempty"`
}

// getDeckCards fetches cards from the database and converts them to both AI and human-readable formats.
func getDeckCards(
	deck map[int]int,
	deckType string,
) ([]*cdb.CardInfoForAI, []*cdb.CardInfoForHuman, error) {
	if len(deck) == 0 {
		return nil, nil, nil
	}

	// Convert map keys to []uint64 for GetCardsByIDs
	ids := make([]uint64, 0, len(deck))
	for code := range deck {
		ids = append(ids, uint64(code))
	}

	// Fetch cards from database
	cards, err := db.GetCardsByIDs(ids)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get %s deck cards: %v", deckType, err)
	}

	// Convert to AI and human formats
	aiCards := make([]*cdb.CardInfoForAI, 0, len(deck))
	humanCards := make([]*cdb.CardInfoForHuman, 0, len(deck))
	for code, count := range deck {
		card, ok := cards[uint64(code)]
		if ok {
			ai := card.ToCardInfoForAI()
			ai.Count = count
			aiCards = append(aiCards, ai)

			// Add deck info for human-readable output
			humanCard := card
			humanCard.Count = optional.Some(count)
			humanCard.Deck = optional.Some(deckType)
			humanCards = append(humanCards, humanCard)
		}
	}
	return aiCards, humanCards, nil
}

// GetCardsByYDKCmd represents the get-cards-by-ydk command
var GetCardsByYDKCmd = &cobra.Command{
	Use:   "get-cards-by-ydk",
	Short: "Get cards from a YDK (Yugioh Deck) file",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if inputFile == "" {
			log.Fatal("input file is required (use --input or -i)")
		}

		// Read the YDK file
		ydkContent, err := os.ReadFile(inputFile)
		if err != nil {
			log.Fatalf("failed to read input file: %v", err)
		}

		// Parse YDK file to get main, extra, side decks
		mainDeck, extraDeck, sideDeck := ydk.ParseYDKFile(string(ydkContent))

		// Fetch and convert deck cards
		mainAI, mainHuman, err := getDeckCards(mainDeck, "main")
		if err != nil {
			log.Fatal(err)
		}

		extraAI, extraHuman, err := getDeckCards(extraDeck, "extra")
		if err != nil {
			log.Fatal(err)
		}

		sideAI, sideHuman, err := getDeckCards(sideDeck, "side")
		if err != nil {
			log.Fatal(err)
		}

		result := &getCardsByYDKOutput{
			Main:  mainAI,
			Extra: extraAI,
			Side:  sideAI,
		}

		// Combine all human-readable cards for CSV
		allHumanCards := append(append(mainHuman, extraHuman...), sideHuman...)

		// Format output
		var output []byte
		switch outputFormat {
		case "json":
			output, err = json.Marshal(result)
			if err != nil {
				log.Fatalf("failed to marshal output: %v", err)
			}
		case "yaml":
			output, err = yaml.Marshal(result)
			if err != nil {
				log.Fatalf("failed to marshal output: %v", err)
			}
		case "csv":
			output, err = formatCardsForHumanToCSV(allHumanCards)
			if err != nil {
				log.Fatalf("failed to format CSV: %v", err)
			}
		default:
			log.Fatalf("unsupported output format: %s (supported: json, yaml, csv)", outputFormat)
		}

		// Write output to file or stdout
		if outputFile != "" {
			err = os.WriteFile(outputFile, output, 0644)
			if err != nil {
				log.Fatalf("failed to write output file: %v", err)
			}
		} else {
			fmt.Println(string(output))
		}
	},
}

func init() {
	GetCardsByYDKCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input YDK file path (required)")
	GetCardsByYDKCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (optional, prints to stdout if not specified)")
	GetCardsByYDKCmd.Flags().StringVar(&outputFormat, "format", "csv", "Output format: json, yaml, csv")
}
