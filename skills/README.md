# Yu-Gi-Oh! AI Tools - Skills

A collection of AI skills for analyzing Yu-Gi-Oh! decks using Claude Code.

## Currently Available Skills

### Core Analysis Skills

| Skill | Description |
|-------|-------------|
| **ygo-deck-analyze** | Main orchestrator that coordinates the complete deck analysis pipeline. Accepts `.ydk` files and generates human-readable reports. |
| **deck-parser** | Parses raw card data to identify engines, starters, extenders, hand traps, and calculates opening probabilities. Performs fact extraction only, no gameplay judgment. |
| **deck-analyzer** | Performs comprehensive gameplay analysis including combo diversity, first-turn endboard quality, weaknesses, handtrap sensitivity, chokepoints, going-second capability, and resource sustainability. Outputs 6-dimensional scores. |

### Utility Skills

| Skill | Description |
|-------|-------------|
| **starter-probability** | Calculates the probability of drawing specific cards (starters, hand traps) in a 5-card starting hand using hypergeometric distribution. |
| **count-deck** | Counts cards in main deck, extra deck, and side deck from a YDK file. |
| **simple-calculator** | Basic arithmetic calculations (+, -, *, /). |

## Analysis Pipeline

The `ygo-deck-analyze` skill orchestrates a 4-step process:

```
.ydk file
    ↓
Step 1: Get raw card data → deck_raw.json
    ↓
Step 2: Parse deck structure → deck_parsed.json
    ↓
Step 3: Gameplay analysis → deck_analysis.json
    ↓
Step 4: Generate report → report.md
```

## Output Format

Each analysis generates:
- `deck_raw.json` - Raw MCP card data
- `deck_parsed.json` - Structured deck information (engines, starters, extenders, hand traps)
- `deck_analysis.json` - Comprehensive gameplay analysis with 6-dimensional scoring
- `report.md` - Human-readable analysis report

## Todo List

### 1. Add English (TCG) Version of Skills

The current skills are primarily in Chinese. Need to:

- [ ] Translate skill descriptions and prompts from Chinese to English
- [ ] Ensure TCG terminology is used correctly (e.g., "hand trap" vs "手坑", "starter" vs "展开点", "board breaker" vs "解场卡")
- [ ] Update reason field templates to use TCG-standard card names and terms

Note: TCG database support is already covered by existing MCP tools.

### 2. Meta-Analyzer Skill

Create a new skill to analyze the competitive meta:

- [ ] Download meta deck lists from online sources (Master Duel, YGOProDeck, etc.)
- [ ] Calculate meta share statistics by archetype
- [ ] Identify common tech cards and hand trap ratios
- [ ] Track card usage trends over time
- [ ] Generate meta reports with deck archetype breakdown

### 3. Combo-Analyzer Skill (Long-term)

An advanced skill using ygopro-core to simulate combos:

- [ ] Integrate ygopro-core Lua engine
- [ ] Build combo tree/graph from deck lists
- [ ] Identify optimal lines and chokepoints
- [ ] Calculate turn 1 board probabilities
- [ ] Find alternative routes when interrupted
- [ ] Generate visual combo maps

## Notes

- All analysis skills are designed for OCG format by default
- The `deck-analyzer` skill requires `resources/handtraps.csv` and `resources/boardbreakers.csv` for complete analysis
- Each judgment in the analysis includes a `reason` field to support reproducible analysis

## License

Part of the yu-gi-oh-ai-tools project.
