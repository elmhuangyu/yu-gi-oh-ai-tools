---
name: count-deck
description: This counts the number of cards in a Yu-Gi-Oh! deck from a YDK file. It parses the file to count cards in the main deck, extra deck, and side deck sections.
---

# Count Deck

This skill counts the number of cards in each deck section from a Yu-Gi-Oh! YDK file.

## Usage

When a user provides a path to a YDK file, use the `scripts/count_deck.py` script to count the cards in each section.

### Example

User: "Count the cards in deck.ydk"
Skill: Calls `python scripts/count_deck.py deck.ydk` and returns:
```
main: 40
extra: 15
side: 15
```
