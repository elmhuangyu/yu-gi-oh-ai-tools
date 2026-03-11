#!/usr/bin/env python3

import sys
import re

def calculate(expression):
    # Remove all whitespace
    expression = expression.replace(" ", "")

    # Check for invalid characters (only numbers and +, -, *, / are allowed)
    if not re.fullmatch(r"[0-9+\-*/.]+", expression):
        print("Error: Invalid characters or parentheses detected. Only numbers and +, -, *, / are supported.")
        sys.exit(1)

    try:
        # Using eval safely by limiting the globals and locals
        # and having already validated the characters in the expression.
        result = eval(expression, {"__builtins__": None}, {})
        return result
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python calculate.py '<expression>'")
        sys.exit(1)
    
    expr = sys.argv[1]
    res = calculate(expr)
    print(res)
