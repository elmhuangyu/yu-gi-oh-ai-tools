# Agent Guidelines

## Project Overview

A Go-based MCP (Model Context Protocol) server that helps locate Yu-Gi-Oh! cards using mycard/ygopro-database.

## Code Style

- Optional Types: A missing value must be distinguished from a valid result, you should use the Optional type to avoid the ambiguity of "magic numbers" (like -1). To implement this, use the [`go-optional`](https://github.com/moznion/go-optional) . See [.ai-docs/optional.md](.ai-docs/go-optional.md) for how to use optional types.

### Testing

See [.ai-docs/testing-code-style.md](.ai-docs/testing-code-style.md).
