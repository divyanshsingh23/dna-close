# G25 Ancestry Estimator

A Go command-line tool that estimates ancestry proportions from G25 coordinates using constrained optimization.

## Overview

This tool allows users to estimate their genetic ancestry by comparing their G25 coordinates against reference populations. It uses a constrained optimization approach to find the best-fitting mixture of reference populations that approximates the target G25 coordinates.

## Features

- Loads reference populations from categorized folders (by time period)
- Parses G25 coordinate files (.txt format)
- Solves a constrained optimization problem to find ancestry proportions
- Outputs clear ancestry estimation results

## Requirements

- Go 1.16+ (no external dependencies needed)

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/yourusername/g25-ancestry-estimator.git
   cd g25-ancestry-estimator
   ```

2. Build the executable:
   ```
   go build -o g25estimator .
   ```

## File Structure

The tool expects reference populations to be organized in the following directory structure:

```
data/
  ├── ancient/
  │   ├── reference1.txt
  │   └── reference2.txt
  ├── medieval/
  │   ├── reference3.txt
  │   └── reference4.txt
  └── modern/
      ├── reference5.txt
      └── reference6.txt
```

## G25 File Format

The tool expects G25 files to be formatted with one sample per line, where each line contains:
- The sample/population label
- 25 G25 coordinates (comma or space separated)

Example:
```
Sample1 0.0124,0.0234,0.0012,0.0055,0.0067,0.0089,0.0012,0.0034,...
Sample2 0.0126 0.0231 0.0015 0.0057 0.0065 0.0084 0.0018 0.0036,...
```

## Usage

1. Prepare your reference population files in the `data/` directory, organized by time period.

2. Run the tool:
   ```
   ./g25estimator
   ```

3. When prompted, enter the path to your target G25 file:
   ```
   Please enter the path to your target G25 file:
   /path/to/your/target.txt
   ```

4. The tool will calculate and display your ancestry proportions.

## How It Works

1. The tool loads all reference populations from the `data/` directory.
2. It parses your target G25 file.
3. It solves an optimization problem:
   - Variables: Weights (ancestry proportions) for each reference population
   - Constraints: All weights ≥ 0 and sum to 1
   - Objective: Minimize Euclidean distance between target vector and weighted average of reference vectors
4. The resulting weights represent estimated ancestry proportions.

## Example Output

```
==================================================
 Ancestry Estimator
==================================================
This tool estimates ancestry proportions from  coordinates.

Loading reference populations...

Reference populations loaded by period:
- ancient: 142 populations
- medieval: 95 populations
- modern: 231 populations

Total reference populations: 468

Please enter the path to your target G25 file:
./mydata.txt

Loading target file: ./mydata.txt

Found 1 sample(s) in target file.

==================================================
Results for: MySample
==================================================
Ancestry Proportions:
--------------------------------------------------
ItalianTuscany                 :  38.25%
SpanishGalicia                 :  24.17%
GreekThessaly                  :  15.88%
BritishWales                   :   9.42%
SerbianBosnia                  :   7.32%
SwedishGotland                 :   4.96%
--------------------------------------------------
```

## License

MIT License
