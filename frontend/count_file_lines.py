#!/usr/bin/env python3
"""Quick discovery script - count lines in a file.
Knuth principle: Make it work, keep it simple for discovery phase."""
from pathlib import Path

# Get file path
filepath = Path(input("Enter file path: ").strip())

# Simple validation - fail fast with clear message
if not filepath.exists():
    print(f"Error: File does not exist: {filepath}")
    exit(1)

if not filepath.is_file():
    print(f"Error: Not a file: {filepath}")
    exit(1)

# Count lines - simple and direct
line_count = 0
try:
    with open(filepath, 'r') as f:
        line_count = sum(1 for _ in f)
except Exception as e:
    print(f"Error reading file: {e}")
    exit(1)

# Output results
print(f"File: {filepath.name}")
print(f"Lines: {line_count}")
print(f"Size: {filepath.stat().st_size} bytes")
