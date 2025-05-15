package main

import (
	"math"
)

// euclideanDistance calculates the Euclidean distance between two vectors
func euclideanDistance(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		panic("vectors must have the same length")
	}
	
	sumSquares := 0.0
	for i := 0; i < len(v1); i++ {
		diff := v1[i] - v2[i]
		sumSquares += diff * diff
	}
	
	return math.Sqrt(sumSquares)
}

// vectorSum adds two vectors element-wise
func vectorSum(v1, v2 []float64) []float64 {
	if len(v1) != len(v2) {
		panic("vectors must have the same length")
	}
	
	result := make([]float64, len(v1))
	for i := 0; i < len(v1); i++ {
		result[i] = v1[i] + v2[i]
	}
	
	return result
}

// vectorSubtract subtracts v2 from v1 element-wise
func vectorSubtract(v1, v2 []float64) []float64 {
	if len(v1) != len(v2) {
		panic("vectors must have the same length")
	}
	
	result := make([]float64, len(v1))
	for i := 0; i < len(v1); i++ {
		result[i] = v1[i] - v2[i]
	}
	
	return result
}

// vectorScale multiplies a vector by a scalar
func vectorScale(v []float64, scalar float64) []float64 {
	result := make([]float64, len(v))
	for i := 0; i < len(v); i++ {
		result[i] = v[i] * scalar
	}
	
	return result
}

// dotProduct calculates the dot product of two vectors
func dotProduct(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		panic("vectors must have the same length")
	}
	
	sum := 0.0
	for i := 0; i < len(v1); i++ {
		sum += v1[i] * v2[i]
	}
	
	return sum
}

// vectorNorm calculates the L2 (Euclidean) norm of a vector
func vectorNorm(v []float64) float64 {
	sumSquares := 0.0
	for _, val := range v {
		sumSquares += val * val
	}
	
	return math.Sqrt(sumSquares)
}

// normalizeVector normalizes a vector to have unit length
func normalizeVector(v []float64) []float64 {
	norm := vectorNorm(v)
	
	// Avoid division by zero
	if norm < 1e-10 {
		return vectorScale(v, 0)
	}
	
	return vectorScale(v, 1/norm)
}

// meanVector calculates the mean vector from a slice of vectors
func meanVector(vectors [][]float64) []float64 {
	if len(vectors) == 0 {
		return nil
	}
	
	dimensions := len(vectors[0])
	result := make([]float64, dimensions)
	
	for _, v := range vectors {
		if len(v) != dimensions {
			panic("all vectors must have the same length")
		}
		
		for i := 0; i < dimensions; i++ {
			result[i] += v[i]
		}
	}
	
	// Divide by the number of vectors
	for i := range result {
		result[i] /= float64(len(vectors))
	}
	
	return result
}
