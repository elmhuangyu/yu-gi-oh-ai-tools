package commands

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/cdb"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

var (
	outputFormat string

	db *cdb.DB
)

// SetDB sets the database instance for commands
func SetDB(database *cdb.DB) {
	db = database
}

func init() {
	GetCardByIDCmd.Flags().StringVar(&outputFormat, "format", "yaml", "Output format: json, yaml, csv")
}

var GetCardByIDCmd = &cobra.Command{
	Use:   "get-card-by-id",
	Short: "Get a card by its ID (passcode)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var id uint64
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
			log.Fatalf("invalid card id: %v", err)
		}

		card, err := db.GetCardByID(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Card with id %d not found", id)
				os.Exit(0)
			}
			log.Fatalf("failed to get card: %v", err)
		}

		cardForAI := card.ToCardInfoForAI()

		var output []byte
		switch outputFormat {
		case "json":
			var err error
			output, err = json.Marshal(cardForAI)
			if err != nil {
				log.Fatalf("failed to marshal card: %v", err)
			}
		case "yaml":
			var err error
			output, err = yaml.Marshal(cardForAI)
			if err != nil {
				log.Fatalf("failed to marshal card: %v", err)
			}
		case "csv":
			headers, rows := cdb.CardInfoForHumanToCSV([]*cdb.CardInfoForHuman{card})
			if headers == nil {
				output = []byte{}
				break
			}
			// Build CSV output using encoding/csv
			var buf strings.Builder
			writer := csv.NewWriter(&buf)
			err := writer.Write(headers)
			if err != nil {
				log.Fatalf("failed to write CSV header: %v", err)
			}
			err = writer.WriteAll(rows)
			if err != nil {
				log.Fatalf("failed to write CSV row: %v", err)
			}
			writer.Flush()
			output = []byte(buf.String())

		default:
			log.Fatalf("unsupported output format: %s (supported: json, yaml)", outputFormat)
		}
		fmt.Println(string(output))
	},
}
