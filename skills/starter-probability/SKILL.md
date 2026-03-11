---
name: starter-probability
description: This calculates the probability of drawing specific cards in a Yu-Gi-Oh! starting hand. Use this when users want to know the odds of opening with certain cards like starters, hand traps.
---

# Starter Probability

This skill calculates the probability of drawing specific cards in a 5-card starting hand using hypergeometric distribution.

## Usage

When a user asks about drawing probabilities, use the `scripts/starter_probability.py` script with the following arguments:

- `--deck-size`: Total number of cards in the main deck (default: 40)
- `--target-count`: Number of target cards in the deck (e.g., number of starters)
- `--hand-size`: Number of cards in the starting hand (default: 5)

## Output Format

The script returns results in YAML format with probabilities for drawing 0, 1, 2, 3, 4, or 5 of the target cards.

### Example 1: Starter cards probability

User: "My deck has 40 cards with 10 starter cards. What's the probability of opening with 1, 2, or 3 starters?"
Skill: Calls `python scripts/starter_probability.py --deck-size 40 --target-count 10`

Returns:
```yaml
deck_size: 40
target_count: 10
hand_size: 5
probabilities:
  0: 21.66%
  1: 41.65%
  2: 27.77%
  3: 7.93%
  4: 0.96%
  5: 0.04%
at_least_1: 78.34%
at_least_2: 36.69%
```

### Example 2: Hand traps probability

User: "40-card deck with 3 Ash Blossom and 3 Maxx 'C'. What's the chance of opening at least 1 hand trap?"
Skill: Calculate combined target count (3 + 3 = 6), then:
`python scripts/starter_probability.py --deck-size 40 --target-count 6`

Returns probability of drawing at least one hand trap.
