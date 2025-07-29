package main

const (
	diffCoeff     = 0.1 // m²/s (smoke diffusivity in air)
	sourceRate    = 50  // kg/m³/s (smoke injection rate)
	density       = 1.0 // \rho = 1 for concentration conservation
	numIterations = 20
	tolerance     = 0.001
)

type Simulation struct {
	mesh      *Mesh
	equations *EquationSystem
	solver    Solver
	time      float64
}

func NewSimulation(canvasWidth, canvasHeight, nx, ny int, scale float64, bcHandler BoundaryHandler) *Simulation {
	mesh := NewMesh(canvasWidth, canvasHeight, nx, ny, scale)
	equationSystem := NewEquationSystem(mesh, diffCoeff, density, bcHandler)
	solver := NewGaussSeidelSolver(numIterations, tolerance)

	return &Simulation{
		mesh:      mesh,
		equations: equationSystem,
		solver:    solver,
		time:      0.0,
	}
}

func (s *Simulation) Step(dt float64) {
	s.equations.AssembleSystem(dt)
	solutions := s.solver.Solve(s.equations.matrix, s.equations.rhs)
	s.mesh.UpdateFromSolutions(solutions)

	s.applySources(dt)
	s.time += dt
}

func (s *Simulation) AddSource(x, y int) {
	s.mesh.SetSource(y, x)
}

func (s *Simulation) applySources(dt float64) {
	for i := 0; i < s.mesh.ny; i++ {
		for j := 0; j < s.mesh.nx; j++ {
			if s.mesh.sources[i][j] {
				s.mesh.phi[i][j] += sourceRate * dt
			}
		}
	}
}
