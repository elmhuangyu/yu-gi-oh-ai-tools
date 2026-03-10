# Agent Guidelines

## Project Overview

A Go-based MCP (Model Context Protocol) server that helps locate Yu-Gi-Oh! cards using mycard/ygopro-database.

## Database

For information about the YGOPRO/YGOCore CDB (SQLite3) database schema, see [.ai-docs/ygopro-cdb.md](.ai-docs/ygopro-cdb.md).

We use [`modernc.org/sqlite`](https://modernc.org/sqlite) (a pure-Go SQLite driver) to access the database.

## Code Style

- Optional Types: A missing value must be distinguished from a valid result, you should use the Optional type to avoid the ambiguity of "magic numbers" (like -1). To implement this, use the [`go-optional`](https://github.com/moznion/go-optional) . See [.ai-docs/optional.md](.ai-docs/go-optional.md) for how to use optional types.

### Testing

See [.ai-docs/testing-code-style.md](.ai-docs/testing-code-style.md).
