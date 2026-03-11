"""Tests for count-deck skill."""

import sys
from pathlib import Path

import pytest

# Add scripts directory to path
sys.path.insert(0, str(Path(__file__).parent.parent / "count-deck" / "scripts"))
from count_deck import count_deck, YDKError


FIXTURES_DIR = Path(__file__).parent / "fixtures" / "decks"


def test_full_deck():
  """Test counting a deck with main, extra, and side sections."""
  result = count_deck(str(FIXTURES_DIR / "sample_deck.ydk"))
  assert result == (5, 2, 1)


def test_empty_deck():
  """Test counting an empty deck (sections but no cards)."""
  result = count_deck(str(FIXTURES_DIR / "empty_deck.ydk"))
  assert result == (0, 0, 0)


def test_only_main():
  """Test counting a deck with only main deck cards."""
  result = count_deck(str(FIXTURES_DIR / "only_main.ydk"))
  assert result == (3, 0, 0)


def test_file_not_found():
  """Test handling of non-existent file."""
  with pytest.raises(YDKError) as exc_info:
    count_deck("nonexistent.ydk")
  assert "not found" in str(exc_info.value).lower()


def test_deck_with_comments_and_empty_lines(tmp_path):
  """Test deck with extra comments and empty lines."""
  deck_file = tmp_path / "test_deck.ydk"
  deck_file.write_text(
    """#created by yugiohcard

#main
89631139
46986414

#extra
6983839

!side
14558127
"""
  )
  result = count_deck(str(deck_file))
  assert result == (2, 1, 1)


def test_deck_case_insensitive(tmp_path):
  """Test that section markers are case-insensitive."""
  deck_file = tmp_path / "case_test.ydk"
  deck_file.write_text(
    """#MAIN
89631139
#EXTRA
6983839
!SIDE
14558127
"""
  )
  result = count_deck(str(deck_file))
  assert result == (1, 1, 1)


def test_deck_ignores_non_numeric_lines(tmp_path):
  """Test that non-numeric lines in sections are ignored."""
  deck_file = tmp_path / "non_numeric_test.ydk"
  deck_file.write_text(
    """#main
89631139
invalid_card_id
46986414
"""
  )
  result = count_deck(str(deck_file))
  assert result == (2, 0, 0)
