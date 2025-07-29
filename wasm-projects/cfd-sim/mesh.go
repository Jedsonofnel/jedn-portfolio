package main

import "fmt"

type Mesh struct {
	phi     [][]float64
	phi0   [][]float64
	sources [][]bool
	nx, ny  int
	dx, dy  float64
	width   float64
	height  float64
}

// resolution is cells/pixel
// scale is m / pixel
func NewMesh(canvasWidth, canvasHeight, nx, ny int, scale float64) *Mesh {
	width := float64(canvasWidth) * scale   // in metres
	height := float64(canvasHeight) * scale // in metres

	dx := width / float64(nx)  // also in metres
	dy := height / float64(ny) // ditto

	return &Mesh{
		phi:     initializeZero2DSlice(nx, ny),
		phi0:   initializeZero2DSlice(nx, ny),
		sources: initializeFalse2DSlice(nx, ny),
		nx:      nx,
		ny:      ny,
		dx:      dx,
		dy:      dy,
		width:   width,
		height:  height,
	}
}

func (m *Mesh) SetSource(i, j int) {
	if !m.inBounds(i, j) {
		panic(fmt.Sprintf("cell coordinates (%d, %d) out of bounds for %dx%d mesh", i, j, m.nx, m.ny))
	}

	m.sources[i][j] = true
}

func (m *Mesh) UpdateFromSolutions(solutions []float64) {
	if len(solutions) != m.nx*m.ny {
		panic(fmt.Sprintf("solutions array of length %d cannot be wrapped across mesh of dimensions [%d, %d]",
			len(solutions), m.nx, m.ny))
	}

	// save previous timestep
	copy2D(m.phi0, m.phi)

	// wrap solutions into mesh
	for n := range(len(solutions)) {
		i := n / m.nx
		j := n % m.nx
		m.phi[i][j] = solutions[n]
	}
}

func (m *Mesh) CellToIndex(i, j int) int {
	return i * m.nx + j
}

func (m *Mesh) inBounds(i, j int) bool {
	return i >= 0 && i < m.ny && j >= 0 && j < m.nx
}

func initializeZero2DSlice(nx, ny int) [][]float64 {
	if nx < 0 || ny < 0 {
		panic("nx and ny cannot be less than 0")
	}

	matrix := make([][]float64, ny)
	for i := range matrix {
		matrix[i] = make([]float64, nx)
	}

	return matrix
}

func initializeFalse2DSlice(nx, ny int) [][]bool {
	if nx < 0 || ny < 0 {
		return nil
	}

	matrix := make([][]bool, ny)
	for i := range matrix {
		matrix[i] = make([]bool, nx)
	}

	return matrix
}

func copy2D(dest, src [][]float64) {
    for i := range src {
        copy(dest[i], src[i])
    }
}
