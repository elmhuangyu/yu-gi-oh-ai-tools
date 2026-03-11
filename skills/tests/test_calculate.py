"""Tests for simple-calculator skill."""

import sys
from pathlib import Path

import pytest

# Add scripts directory to path
sys.path.insert(0, str(Path(__file__).parent.parent / "simple-calculator" / "scripts"))
from calculate import calculate, CalculateError


def test_basic_addition():
  """Test simple addition."""
  assert calculate("2+2") == 4
  assert calculate("10+5") == 15


def test_basic_subtraction():
  """Test simple subtraction."""
  assert calculate("5-3") == 2
  assert calculate("10-20") == -10


def test_basic_multiplication():
  """Test simple multiplication."""
  assert calculate("3*4") == 12
  assert calculate("7*8") == 56


def test_basic_division():
  """Test simple division."""
  assert calculate("10/2") == 5
  assert calculate("7/2") == 3.5


def test_complex_expression():
  """Test combined operations."""
  assert calculate("2+3*4") == 14
  assert calculate("10+20/5-3") == 11


def test_float_numbers():
  """Test floating point numbers."""
  assert calculate("3.5+2.5") == 6.0
  assert calculate("10.5/2") == 5.25


def test_whitespace_handling():
  """Test that whitespace is properly removed."""
  assert calculate("2 + 2") == 4
  assert calculate(" 10   +   20 ") == 30
  assert calculate("3 * 4 + 2") == 14


def test_invalid_characters():
  """Test that invalid characters are rejected."""
  invalid_cases = [
    "2+2@",  # @ symbol
    "2^10",  # exponent operator
    "(1+2)*3",  # parentheses
    "2&3",  # bitwise AND
  ]

  for expr in invalid_cases:
    with pytest.raises(CalculateError):
      calculate(expr)


def test_parentheses_rejected():
  """Test that parentheses are rejected (per current implementation)."""
  with pytest.raises(CalculateError) as exc_info:
    calculate("(2+3)*4")
  assert "parentheses" in str(exc_info.value).lower()


def test_invalid_expression():
  """Test handling of invalid expressions."""
  with pytest.raises(CalculateError):
    calculate("2+/3")


def test_empty_expression():
  """Test handling of empty expression."""
  with pytest.raises(CalculateError):
    calculate("")


def test_negative_numbers():
  """Test handling of negative numbers."""
  assert calculate("-5+3") == -2
  assert calculate("10*-2") == -20
