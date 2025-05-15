package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	dataDir = "./data"
)

func main() {
	// Welcome message
	fmt.Println("==================================================")
	fmt.Println("G25 Ancestry Estimator")
	fmt.Println("==================================================")
	fmt.Println("This tool estimates ancestry proportions from G25 coordinates.")

	// Load reference populations
	fmt.Println("\nLoading reference populations...")
	populations, err := loadReferencePopulations(dataDir)
	if err != nil {
		fmt.Printf("Error loading reference populations: %v\n", err)
		os.Exit(1)
	}

	// Print summary of loaded data
	var totalPopulations int
	fmt.Println("\nReference populations loaded by period:")
	for period, pops := range populations {
		fmt.Printf("- %s: %d populations\n", period, len(pops))
		totalPopulations += len(pops)
	}
	fmt.Printf("\nTotal reference populations: %d\n", totalPopulations)

	// Get target file from user
	fmt.Println("\nPlease enter the path to your target G25 file:")
	reader := bufio.NewReader(os.Stdin)
	targetPath, _ := reader.ReadString('\n')
	targetPath = strings.TrimSpace(targetPath)

	// Load target coordinates
	fmt.Printf("\nLoading target file: %s\n", targetPath)
	targetSamples, err := parseG25File(targetPath)
	if err != nil {
		fmt.Printf("Error loading target file: %v\n", err)
		os.Exit(1)
	}

	// Run ancestry estimation for each target sample
	fmt.Printf("\nFound %d sample(s) in target file.\n", len(targetSamples))
	
	// Flatten all populations for the solver
	allPops := flattenPopulations(populations)
	
	for _, sample := range targetSamples {
		fmt.Printf("\n==================================================\n")
		fmt.Printf("Results for: %s\n", sample.Label)
		fmt.Printf("==================================================\n")
		
		// Solve for ancestry proportions
		results, err := solveAncestryProportions(sample.Coordinates, allPops)
		if err != nil {
			fmt.Printf("Error solving ancestry proportions: %v\n", err)
			continue
		}
		
		// Display results
		displayResults(results)
	}
}

// loadReferencePopulations loads all reference populations from the data directory
func loadReferencePopulations(baseDir string) (map[string][]Population, error) {
	result := make(map[string][]Population)
	
	// Read the base directory
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}
	
	// Process each subdirectory (period)
	for _, entry := range entries {
		if entry.IsDir() {
			periodName := entry.Name()
			periodPath := filepath.Join(baseDir, periodName)
			
			// Read all files in this period directory
			files, err := os.ReadDir(periodPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read period directory %s: %w", periodName, err)
			}
			
			// Process each reference file
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
					filePath := filepath.Join(periodPath, file.Name())
					populations, err := parseG25File(filePath)
					if err != nil {
						return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
					}
					
					// Add to the appropriate period
					result[periodName] = append(result[periodName], populations...)
				}
			}
		}
	}
	
	return result, nil
}

// flattenPopulations converts the period-organized map into a single slice of populations
func flattenPopulations(populations map[string][]Population) []Population {
	var result []Population
	
	for _, pops := range populations {
		result = append(result, pops...)
	}
	
	return result
}

// displayResults prints the ancestry proportions in a readable format
func displayResults(results []AncestryResult) {
	// Sort results by proportion (descending)
	sortAncestryResults(results)
	
	// Print the top results
	fmt.Println("Ancestry Proportions:")
	fmt.Println("--------------------------------------------------")
	
	// Calculate total to ensure it sums to 100%
	var total float64
	for _, result := range results {
		if result.Proportion > 0.001 { // Only count non-trivial contributions
			total += result.Proportion
		}
	}
	
	// Print each contribution
	for _, result := range results {
		// Only show populations with meaningful contributions (>0.1%)
		if result.Proportion > 0.001 {
			// Normalize to ensure sum is exactly 100%
			normalizedProp := (result.Proportion / total) * 100
			fmt.Printf("%-30s: %6.2f%%\n", result.Population.Label, normalizedProp)
		}
	}
	
	fmt.Println("--------------------------------------------------")
}
