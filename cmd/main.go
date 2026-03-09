package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/elmhuangyu/yu-gi-oh-mcp/lib/cdb"
	"github.com/elmhuangyu/yu-gi-oh-mcp/lib/git"
)

func getLangFromEnvOrDefault() string {
	if os.Getenv("APP_LANG") != "" {
		return os.Getenv("APP_LANG")
	}
	return cdb.SupportLanguages[0]
}

var (
	flagLang = flag.String("lang", getLangFromEnvOrDefault(), fmt.Sprintf("language of this MCP server, currently only support %v", cdb.SupportLanguages))
)

const (
	repoClonePath = "/tmp/yugioh-cdb"
	cdbRepoURL    = "https://github.com/mycard/ygopro-database.git"
)

func main() {
	flag.Parse()

	if !slices.Contains(cdb.SupportLanguages, *flagLang) {
		log.Fatalf("lang is not supported: %s", *flagLang)
	}

	repo := git.NewRepo(repoClonePath, cdbRepoURL)
	if err := repo.EnsureRepoUpToDate(); err != nil {
		log.Fatalf("cdb repo failed to update: %v", err)
	}

}
