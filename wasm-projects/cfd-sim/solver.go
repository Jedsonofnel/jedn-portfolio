package main

import (
	"fmt"
	"math"
)

const DEBUG = false

type Solver interface {
	Solve(matrix [][]float64, rhs []float64) []float64
}

type GaussSeidelSolver struct {
	maxIterations int
	tolerance     float64
}

func NewGaussSeidelSolver(maxIterations int, tolerance float64) *GaussSeidelSolver {
	return &GaussSeidelSolver{
		maxIterations: maxIterations,
		tolerance:     tolerance,
	}
}

func (gs *GaussSeidelSolver) Solve(matrix [][]float64, rhs []float64) []float64 {
	n := len(rhs)
	x := make([]float64, n)

	for i := range x {
		x[i] = 0.0
	}

	for iter := 0; iter < gs.maxIterations; iter++ {
		residual := 0.0

		for i := 0; i < n; i++ {
			sum := 0.0

			for j := 0; j < n; j++ {
				if i != j {
					sum += matrix[i][j] * x[j]
				}
			}

			newValue := (rhs[i] - sum) / matrix[i][i]
			residual += math.Abs(newValue - x[i])
			x[i] = newValue
		}

		if residual < gs.tolerance {
			if DEBUG {
				fmt.Printf("Gauss-Seidel converged after %d iterations\n", iter)
			}
			break
		}

		if iter%10 == 0 {
			if DEBUG {
				fmt.Printf("Iteration %d, residual: %.6f\n", iter, residual)
			}
		}
	}

	return x
}

type JacobiSolver struct {
	maxIterations int
	tolerance     float64
}

func NewJacobiSolver(maxIterations int, tolerance float64) *JacobiSolver {
	return &JacobiSolver{
		maxIterations: maxIterations,
		tolerance:     tolerance,
	}
}

func (j *JacobiSolver) Solve(matrix [][]float64, rhs []float64) []float64 {
	n := len(rhs)
	x := make([]float64, n)
	xNew := make([]float64, n)

	for iter := 0; iter < j.maxIterations; iter++ {
		for i := 0; i < n; i++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				if i != k {
					sum += matrix[i][k] * x[k]
				}
			}
			xNew[i] = (rhs[i] - sum) / matrix[i][i]
		}

		// Copy xNew to x
		copy(x, xNew)
	}

	return x
}
