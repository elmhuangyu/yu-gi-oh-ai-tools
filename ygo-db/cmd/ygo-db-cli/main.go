package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/cmd/ygo-db-cli/commands"
	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/cdb"
	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/git"
	"github.com/spf13/cobra"
)

const (
	cdbRepoURL = "https://github.com/mycard/ygopro-database.git"
)

var (
	flagLang string
	db       *cdb.DB
	gitRepo  *git.Repo
)

func getBasePath() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "ygo-db")
	}
	return filepath.Join(os.Getenv("HOME"), ".config/ygo-db")
}

var rootCmd = &cobra.Command{
	Use:   "ygo-db-cli",
	Short: "A CLI tool to query Yu-Gi-Oh! card database",
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&flagLang,
		"lang",
		cdb.SupportLanguages[0],
		fmt.Sprintf("Language for card data: %v", cdb.SupportLanguages))

	rootCmd.AddCommand(commands.GetCardByIDCmd)
	rootCmd.AddCommand(commands.GetCardsByNameCmd)
	rootCmd.AddCommand(commands.GetCardsByArchetypesCmd)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Validate lang flag
		validLang := slices.Contains(cdb.SupportLanguages, flagLang)
		if !validLang {
			return fmt.Errorf("invalid language: %s (supported: %v)", flagLang, cdb.SupportLanguages)
		}

		basePath := getBasePath()

		var err error
		gitRepo = git.NewRepo(basePath, cdbRepoURL)
		if err = gitRepo.EnsureRepoUpToDate(); err != nil {
			log.Printf("cdb repo failed to update: %v", err)
		}

		db, err = cdb.New(gitRepo, basePath, flagLang, false)
		if err != nil {
			return err
		}

		// Set db instance for commands
		commands.SetDB(db)

		return nil
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("command execution failed: %v", err)
	}
}
