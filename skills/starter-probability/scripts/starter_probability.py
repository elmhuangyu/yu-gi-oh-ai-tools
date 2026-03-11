#!/usr/bin/env python3
"""Calculate Yu-Gi-Oh! starting hand probabilities using hypergeometric distribution."""

import argparse
from math import comb
from typing import Dict


def hypergeometric_pmf(
  population_size: int, success_states: int, draws: int, observed_successes: int
) -> float:
  """Calculate probability mass function for hypergeometric distribution.

  P(X = k) = C(K, k) * C(N-K, n-k) / C(N, n)

  Args:
      population_size: Total population size (N) - deck size
      success_states: Number of success states in population (K) - target cards
      draws: Number of draws (n) - hand size
      observed_successes: Number of observed successes (k) - target cards drawn

  Returns:
      Probability of drawing exactly observed_successes target cards
  """
  if observed_successes < 0:
    return 0.0
  if observed_successes > min(success_states, draws):
    return 0.0
  if observed_successes < max(0, draws - (population_size - success_states)):
    return 0.0

  numerator = comb(success_states, observed_successes) * comb(
    population_size - success_states, draws - observed_successes
  )
  denominator = comb(population_size, draws)
  return numerator / denominator


def calculate_probabilities(
  deck_size: int, target_count: int, hand_size: int = 5
) -> Dict[str, float]:
  """Calculate probabilities for drawing 0 to hand_size target cards.

  Args:
      deck_size: Total number of cards in the main deck
      target_count: Number of target cards in the deck
      hand_size: Number of cards in starting hand (default: 5)

  Returns:
      Dictionary with probability results
  """
  results = {
    "deck_size": deck_size,
    "target_count": target_count,
    "hand_size": hand_size,
    "probabilities": {},
  }

  # Calculate individual probabilities
  for k in range(hand_size + 1):
    prob = hypergeometric_pmf(deck_size, target_count, hand_size, k)
    results["probabilities"][str(k)] = round(prob * 100, 2)

  # Calculate cumulative probabilities
  at_least_1 = sum(
    hypergeometric_pmf(deck_size, target_count, hand_size, k) for k in range(1, hand_size + 1)
  )
  results["at_least_1"] = round(at_least_1 * 100, 2)

  at_least_2 = sum(
    hypergeometric_pmf(deck_size, target_count, hand_size, k) for k in range(2, hand_size + 1)
  )
  results["at_least_2"] = round(at_least_2 * 100, 2)

  return results


def to_yaml(data: dict, indent: int = 0) -> str:
  """Convert dictionary to YAML format string."""
  lines = []
  prefix = "  " * indent

  for key, value in data.items():
    if isinstance(value, dict):
      lines.append(f"{prefix}{key}:")
      lines.append(to_yaml(value, indent + 1))
    elif isinstance(value, float):
      lines.append(f"{prefix}{key}: {value}%")
    elif isinstance(value, int):
      lines.append(f"{prefix}{key}: {value}")
    else:
      lines.append(f"{prefix}{key}: {value}")

  return "\n".join(lines)


def main():
  parser = argparse.ArgumentParser(description="Calculate Yu-Gi-Oh! starting hand probabilities")
  parser.add_argument(
    "--deck-size", "-d", type=int, default=40, help="Main deck size (default: 40)"
  )
  parser.add_argument(
    "--target-count", "-t", type=int, required=True, help="Number of target cards in deck"
  )
  parser.add_argument("--hand-size", type=int, default=5, help="Starting hand size (default: 5)")

  args = parser.parse_args()

  # Validate inputs
  if args.deck_size <= 0:
    raise ValueError("Deck size must be positive")
  if args.target_count < 0:
    raise ValueError("Target count cannot be negative")
  if args.target_count > args.deck_size:
    raise ValueError("Target count cannot exceed deck size")
  if args.hand_size <= 0:
    raise ValueError("Hand size must be positive")
  if args.hand_size > args.deck_size:
    raise ValueError("Hand size cannot exceed deck size")

  results = calculate_probabilities(args.deck_size, args.target_count, args.hand_size)
  print(to_yaml(results))


if __name__ == "__main__":
  main()
