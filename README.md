# G25 Coordinate Matcher

A terminal-based application to find the closest matching samples based on G25 coordinates.

## Requirements

- Python 3.7+
- No external dependencies required

## File Requirements

1. `xyz.txt`: Your G25 coordinates (25 comma-separated values on a single line)
2. `samples.txt`: Database of samples with their G25 coordinates
   - One sample per line
   - Format: SampleName,coord1,coord2,...,coord25

## Usage

1. Place your G25 coordinates in `xyz.txt`
2. Place your sample database in `samples.txt`
3. Run the program:
   ```bash
   python main.py
   ```

## Features

- Euclidean distance calculation
- Top 10 closest matches display
- Sample count functionality
- Error handling for missing/invalid files
- Efficient processing for large sample files
