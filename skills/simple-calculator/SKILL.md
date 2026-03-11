---
name: simple-calculator
description: This performs basic arithmetic calculations (+, -, *, /). Use this skill whenever the user asks for a simple calculation involving these operators. It does not support parentheses or advanced math functions.
---

# Simple Calculator

This skill provides a simple way to perform basic arithmetic operations using a Python script.

## Constraints
- Does NOT support parentheses `()`.
- Only supports basic operators: `+`, `-`, `*`, `/`.
- Inputs should be simple arithmetic expressions.

## Usage
When a user provides an arithmetic expression, use the `scripts/calculate.py` script to get the result.

### Example
User: "What is 10 + 5 * 2?"
Skill: Calls `python scripts/calculate.py "10 + 5 * 2"` and returns the result.
