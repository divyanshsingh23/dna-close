package main

import (
	"fmt"
	"math"
	"sort"
)

// AncestryResult represents a single population contribution to ancestry
type AncestryResult struct {
	Population  Population
	Proportion  float64 // Value between 0 and 1
}

// solveAncestryProportions finds the optimal mixture of reference populations
// that best approximates the target coordinates
func solveAncestryProportions(target []float64, populations []Population) ([]AncestryResult, error) {
	if len(target) != 25 {
		return nil, fmt.Errorf("target must have exactly 25 dimensions, got %d", len(target))
	}
	
	if len(populations) == 0 {
		return nil, fmt.Errorf("no reference populations provided")
	}
	
	// Check that all populations have 25 dimensions
	for i, pop := range populations {
		if len(pop.Coordinates) != 25 {
			return nil, fmt.Errorf("population %d (%s) has %d dimensions, expected 25", 
				i, pop.Label, len(pop.Coordinates))
		}
	}
	
	// Number of populations
	n := len(populations)
	
	// Initialize weights equally
	weights := make([]float64, n)
	for i := range weights {
		weights[i] = 1.0 / float64(n)
	}
	
	// Perform projected gradient descent
	const (
		maxIterations = 10000
		learningRate  = 0.01
		convergenceThreshold = 1e-6
	)
	
	var prevError float64 = math.Inf(1)
	
	for iter := 0; iter < maxIterations; iter++ {
		// Calculate current error
		weighted := weightedAverage(populations, weights)
		currentError := euclideanDistance(target, weighted)
		
		// Check for convergence
		if math.Abs(prevError - currentError) < convergenceThreshold {
			break
		}
		prevError = currentError
		
		// Calculate gradients
		gradients := make([]float64, n)
		for i := range gradients {
			// Partial derivative of squared Euclidean distance with respect to weight i
			grad := 0.0
			for d := 0; d < 25; d++ {
				weightedSum := 0.0
				for j := range weights {
					weightedSum += weights[j] * populations[j].Coordinates[d]
				}
				grad += 2 * (weightedSum - target[d]) * populations[i].Coordinates[d]
			}
			gradients[i] = grad
		}
		
		// Update weights using gradient descent
		for i := range weights {
			weights[i] -= learningRate * gradients[i]
		}
		
		// Project weights back to satisfy constraints (non-negative and sum to 1)
		projectToSimplex(weights)
	}
	
	// Build results
	results := make([]AncestryResult, n)
	for i := range populations {
		results[i] = AncestryResult{
			Population: populations[i],
			Proportion: weights[i],
		}
	}
	
	return results, nil
}

// weightedAverage calculates the weighted average of population coordinates
func weightedAverage(populations []Population, weights []float64) []float64 {
	result := make([]float64, 25)
	
	for d := 0; d < 25; d++ {
		for i, pop := range populations {
			result[d] += weights[i] * pop.Coordinates[d]
		}
	}
	
	return result
}

// projectToSimplex projects the weights vector onto the probability simplex
// (all weights non-negative and sum to 1)
func projectToSimplex(weights []float64) {
	// First, project to non-negative orthant
	for i := range weights {
		if weights[i] < 0 {
			weights[i] = 0
		}
	}
	
	// Then, normalize to sum to 1
	sum := 0.0
	for _, w := range weights {
		sum += w
	}
	
	// If sum is very close to zero, distribute equally
	if sum < 1e-10 {
		for i := range weights {
			weights[i] = 1.0 / float64(len(weights))
		}
		return
	}
	
	// Otherwise normalize
	for i := range weights {
		weights[i] /= sum
	}
}

// sortAncestryResults sorts results in descending order of proportion
func sortAncestryResults(results []AncestryResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Proportion > results[j].Proportion
	})
}
