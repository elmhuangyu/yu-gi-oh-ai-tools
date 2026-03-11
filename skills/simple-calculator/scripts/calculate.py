#!/usr/bin/env python3

import sys
import re


class CalculateError(Exception):
  """Exception raised for calculation errors."""

  pass


def calculate(expression):
  """
  Calculate the result of a mathematical expression.

  Args:
      expression: A string containing a mathematical expression

  Returns:
      The result of the calculation

  Raises:
      CalculateError: If the expression contains invalid characters
                      or cannot be evaluated
  """
  # Remove all whitespace
  expression = expression.replace(" ", "")

  # Check for invalid characters (only numbers and +, -, *, / are allowed)
  if not re.fullmatch(r"[0-9+\-*/.]+", expression):
    raise CalculateError(
      "Invalid characters or parentheses detected. Only numbers and +, -, *, / are supported."
    )

  try:
    # Using eval safely by limiting the globals and locals
    # and having already validated the characters in the expression.
    result = eval(expression, {"__builtins__": None}, {})
    return result
  except Exception as e:
    raise CalculateError(str(e))


def main():
  """CLI entry point for calculate."""
  if len(sys.argv) != 2:
    print("Usage: python calculate.py '<expression>'")
    sys.exit(1)

  expr = sys.argv[1]

  try:
    res = calculate(expr)
    print(res)
  except CalculateError as e:
    print(f"Error: {e}")
    sys.exit(1)


if __name__ == "__main__":
  main()
