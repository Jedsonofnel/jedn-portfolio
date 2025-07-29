package main

// TODO Add velocity field interface: velocityField VelocityField
type EquationSystem struct {
	mesh      *Mesh
	matrix    [][]float64
	rhs       []float64
	diffCoeff float64
	density   float64
	bcHandler BoundaryHandler
}

func NewEquationSystem(mesh *Mesh, diffCoeff, density float64, bcHandler BoundaryHandler) *EquationSystem {
	nCells := mesh.nx * mesh.ny

	return &EquationSystem{
		mesh:      mesh,
		matrix:    initialize2DSystemMatrix(nCells),
		rhs:       make([]float64, nCells),
		diffCoeff: diffCoeff,
		density:   density,
		bcHandler: bcHandler,
	}
}

func (eq *EquationSystem) AssembleSystem(dt float64) {
	eq.clearMatrix()
	diffX := eq.diffCoeff * eq.mesh.dy / eq.mesh.dx
	diffY := eq.diffCoeff * eq.mesh.dx / eq.mesh.dy

	for i := range eq.mesh.ny {
		for j := range eq.mesh.nx {
			cellIndex := eq.mesh.CellToIndex(i, j)

			aE, aW, aN, aS := eq.bcHandler.HandleBoundaryCoefficients(i, j, eq.mesh.nx, eq.mesh.ny, diffX, diffY)

			aP0 := (eq.mesh.dy * eq.mesh.dx * eq.density) / dt
			aP := aE + aW + aN + aS + aP0

			// Set matrix coefficients
			eq.setNeighborCoefficients(cellIndex, i, j, aE, aW, aN, aS)
			eq.matrix[cellIndex][cellIndex] = aP

			// Set RHS
			phi0 := eq.mesh.phi0[i][j]
			eq.rhs[cellIndex] = aP0 * phi0
			eq.rhs[cellIndex] += eq.bcHandler.HandleBoundaryRHS(
				i, j, eq.mesh.nx, eq.mesh.ny, aE, aW, aN, aS, eq.mesh.dx, eq.mesh.dy)
		}
	}
}

func (eq *EquationSystem) calculateCoefficients(i, j int, diff_x, diff_y float64) (a_E, a_W, a_N, a_S float64) {
	// East coefficient
	if j == eq.mesh.nx-1 {
		a_E = 2 * diff_x
	} else {
		a_E = diff_x
	}

	// West coefficient
	if j == 0 {
		a_W = 2 * diff_x
	} else {
		a_W = diff_x
	}

	// North coefficient
	if i == 0 {
		a_N = 2 * diff_y
	} else {
		a_N = diff_y
	}

	// South coefficient
	if i == eq.mesh.ny-1 {
		a_S = 2 * diff_y
	} else {
		a_S = diff_y
	}

	return
}

func (eq *EquationSystem) setNeighborCoefficients(cellIndex, i, j int, a_E, a_W, a_N, a_S float64) {
	// East neighbor
	if j < eq.mesh.nx-1 {
		eastIndex := eq.mesh.CellToIndex(i, j+1)
		eq.matrix[cellIndex][eastIndex] = -a_E
	}

	// West neighbor
	if j > 0 {
		westIndex := eq.mesh.CellToIndex(i, j-1)
		eq.matrix[cellIndex][westIndex] = -a_W
	}

	// North neighbor
	if i > 0 {
		northIndex := eq.mesh.CellToIndex(i-1, j)
		eq.matrix[cellIndex][northIndex] = -a_N
	}

	// South neighbor
	if i < eq.mesh.ny-1 {
		southIndex := eq.mesh.CellToIndex(i+1, j)
		eq.matrix[cellIndex][southIndex] = -a_S
	}
}

func (eq *EquationSystem) clearMatrix() {
	for i := range eq.matrix {
		for j := range eq.matrix[i] {
			eq.matrix[i][j] = 0.0
		}
	}
	for i := range eq.rhs {
		eq.rhs[i] = 0.0
	}
}

func initialize2DSystemMatrix(nCells int) [][]float64 {
	matrix := make([][]float64, nCells)
	for row := range matrix {
		matrix[row] = make([]float64, nCells)
	}
	return matrix
}
