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
	outputFormat string

	db *cdb.DB
)

// SetDB sets the database instance for commands
func SetDB(database *cdb.DB) {
	db = database
}

func init() {
	GetCardByIDCmd.Flags().StringVar(&outputFormat, "format", "yaml", "Output format: json, yaml")
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
		default:
			log.Fatalf("unsupported output format: %s (supported: json, yaml)", outputFormat)
		}
		fmt.Println(string(output))
	},
}
