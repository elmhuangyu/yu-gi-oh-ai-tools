# Yu-Gi-Oh! AI Tools

**English** | [简体中文](README-zh.md)

A specialized toolkit for analyzing Yu-Gi-Oh! decks using AI Agent Skills.

## Overview

This project provides a suite of tools and Agent Skills to help AI assistants (like Gemini CLI or Claude Code) understand and analyze Yu-Gi-Oh! decklists. It automates the process of parsing `.ydk` files, calculating probabilities, and generating professional-grade gameplay evaluations.

## Project Structure

```text
yu-gi-oh-ai-tools/
├── skills/                    # Agent Skills for Gemini CLI / Claude Code
│   ├── ygo-deck-analyze/      # Main orchestration logic
│   ├── deck-parser/           # Card classification & probability engine
│   ├── deck-analyzer/         # Gameplay & board-breaking analysis
│   ├── starter-probability/   # Hypergeometric calculator
│   ├── count-deck/            # Deck size validator
│   └── simple-calculator/     # Utility math
│
├── ygo-db/                    # Go-based CLI card database
│   └── cmd/ygo-db-cli/        # High-performance card lookup
│
└── test-dir-zh-cn/            # Sample decks for testing
```

## Features

### Core Analysis Skills

| Skill | Description |
|-------|-------------|
| **ygo-deck-analyze** | The entry point. Coordinates the full pipeline from raw data to final report. |
| **deck-parser** | Extracts engine cores, starters, extenders, and hand traps. Focuses on factual data extraction. |
| **deck-analyzer** | Evaluates deck performance: combo diversity, endboard strength, choke points, and resource loops. |

### Utility Skills

| Skill | Description |
|-------|-------------|
| **starter-probability** | Uses hypergeometric distribution to calculate the odds of drawing specific card combinations. |
| **count-deck** | Quickly counts main, extra, and side deck sizes from YDK files. |
| **ygoprodeck-get-meta** | Fetches the latest meta decklists from YGOPRODeck. |

### Database Tools

- **ygo-db-cli**: A fast Go-based CLI for querying local card databases by ID, name, archetype, or YDK file.

## Analysis Workflow

The `ygo-deck-analyze` skill automates the following pipeline:

1.  **Data Retrieval**: Fetch raw card metadata via `ygo-db-cli`.
2.  **Structural Parsing**: Identify starters, extenders, and engine pieces.
3.  **Gameplay Analysis**: Evaluate deck metrics (consistency, ceiling, etc.).
4.  **Reporting**: Generate a formatted `report.md`.

## Output Artifacts

Each analysis run produces:
- `deck_raw.csv`: Raw metadata from the database.
- `deck_parsed.json`: Categorized card roles and engine components.
- `deck_analysis.json`: Quantitative metrics and 6-dimensional scoring.
- `report.md`: Human-readable summary and strategic advice.

## Six-Dimensional Scoring

| Metric | Focus |
|-----------|-------------|
| **Consistency** | Reliability of opening hands / starter density. |
| **Ceiling** | Maximum endboard power when uninterrupted. |
| **Resilience** | Ability to play through hand traps or disruption. |
| **Going Second** | Board-breaking potential and comeback ability. |
| **Grind** | Resource recursion and long-game sustainability. |
| **Diversity** | Variety of combo paths and playstyles. |

## Quick Start

```bash
# Get raw data for a deck
ygo-db-cli get-cards-by-ydk -i deck.ydk -o deck_raw.csv

# Run analysis via Gemini CLI / Claude Code
# Use the 'ygo-deck-analyze' skill on the generated CSV
```

## Requirements

- **AI Skills**: Gemini CLI, Claude Code, or compatible AI agent environments.
- **ygo-db**: Go 1.21+
- **Python**: 3.14+ (for utility scripts)

## License

Apache 2.0

---

*Yu-Gi-Oh! is a trademark of Konami Digital Entertainment. This project is not affiliated with, endorsed, sponsored, or specifically approved by Konami Digital Entertainment.*
