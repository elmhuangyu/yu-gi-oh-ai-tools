#!/usr/bin/env python3
"""
Count the number of cards in each deck section of a Yu-Gi-Oh! YDK file.

Usage: python count_deck.py <path_to_ydk_file>
"""

import sys


def count_deck(ydk_path):
  """
  Parse a YDK file and count cards in main, extra, and side decks.

  Args:
      ydk_path: Path to the YDK file

  Returns:
      Tuple of (main_count, extra_count, side_count)
  """
  main_count = 0
  extra_count = 0
  side_count = 0

  current_section = None

  try:
    with open(ydk_path, "r", encoding="utf-8") as f:
      for line in f:
        line = line.strip()

        # Skip empty lines
        if not line:
          continue

        # Check for section markers
        if line.lower() == "#main":
          current_section = "main"
          continue
        elif line.lower() == "#extra":
          current_section = "extra"
          continue
        elif line.lower() == "!side":
          current_section = "side"
          continue

        # Skip other comment lines (starting with #)
        if line.startswith("#"):
          continue

        # If we have a valid section and the line is a card ID (numeric)
        # Count it
        if current_section and line.isdigit():
          if current_section == "main":
            main_count += 1
          elif current_section == "extra":
            extra_count += 1
          elif current_section == "side":
            side_count += 1

  except FileNotFoundError:
    print(f"Error: File '{ydk_path}' not found.")
    sys.exit(1)
  except Exception as e:
    print(f"Error: {e}")
    sys.exit(1)

  return main_count, extra_count, side_count


def main():
  if len(sys.argv) != 2:
    print("Usage: python count_deck.py <path_to_ydk_file>")
    sys.exit(1)

  ydk_path = sys.argv[1]

  main_count, extra_count, side_count = count_deck(ydk_path)

  print(f"main: {main_count}")
  print(f"extra: {extra_count}")
  print(f"side: {side_count}")


if __name__ == "__main__":
  main()
