package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Population represents a reference population with its G25 coordinates
type Population struct {
	Label       string    // The name/label of the population
	Coordinates []float64 // The 25-dimensional G25 vector
	Period      string    // The time period this population belongs to (derived from folder)
}

// parseG25File reads a G25 file and returns a slice of Population structs
func parseG25File(filePath string) ([]Population, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var populations []Population
	period := extractPeriodFromPath(filePath)
	
	scanner := bufio.NewScanner(file)
	lineNum := 0
	
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// The format appears to be: Label,coord1,coord2,...,coord25
		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format, expected 'Label,coord1,coord2,...'", lineNum)
		}
		
		label := parts[0]
		coordsStr := parts[1]
		
		// Split the coordinates part by comma
		coordFields := strings.Split(coordsStr, ",")
		if len(coordFields) != 25 {
			return nil, fmt.Errorf("line %d: expected exactly 25 coordinates, got %d", lineNum, len(coordFields))
		}
		
		// Parse the 25 coordinates
		coordinates := make([]float64, 25)
		for i := 0; i < 25; i++ {
			value, err := strconv.ParseFloat(strings.TrimSpace(coordFields[i]), 64)
			if err != nil {
				return nil, fmt.Errorf("line %d, coordinate %d: invalid float value '%s': %w", 
					lineNum, i+1, coordFields[i], err)
			}
			coordinates[i] = value
		}
		
		// Add to our list of populations
		populations = append(populations, Population{
			Label:       label,
			Coordinates: coordinates,
			Period:      period,
		})
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	
	if len(populations) == 0 {
		return nil, fmt.Errorf("no valid populations found in file")
	}
	
	return populations, nil
}

// extractPeriodFromPath attempts to determine the period from the file path
func extractPeriodFromPath(filePath string) string {
	// Get directory name which should be the period
	dir := filepath.Dir(filePath)
	period := filepath.Base(dir)
	
	// If this is the base data directory, set period to "unknown"
	if period == "data" || period == "." {
		return "unknown"
	}
	
	return period
}
