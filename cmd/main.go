package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/elmhuangyu/yu-gi-oh-mcp/lib/cdb"
	"github.com/elmhuangyu/yu-gi-oh-mcp/lib/git"
	"github.com/elmhuangyu/yu-gi-oh-mcp/lib/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func getLangFromEnvOrDefault() string {
	if os.Getenv("APP_LANG") != "" {
		return os.Getenv("APP_LANG")
	}
	return cdb.SupportLanguages[0]
}

var (
	flagLang = flag.String("lang", getLangFromEnvOrDefault(), fmt.Sprintf("language of this MCP server, currently only support %v", cdb.SupportLanguages))
	flagPort = flag.Int("port", 8000, "port of this MCP server")
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
	log.Println("cdb repo updated")

	db, err := cdb.New(repo, repoClonePath, *flagLang)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "yugioh-card-database", Version: "v1.0.0"}, nil)
	tools.Tools(server, db)

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *flagPort),
		Handler: handler,
	}

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Printf("MCP Server Listening on port: %d...", *flagPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
