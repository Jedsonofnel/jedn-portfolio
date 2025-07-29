package main

type BoundaryType int

const (
	Dirichlet BoundaryType = iota
	Neumann
	Outflow
)

type BoundaryCondition struct {
	Type  BoundaryType
	Value float64
}

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type BoundaryConditions struct {
	North BoundaryCondition `json:"north"`
	East  BoundaryCondition `json:"east"`
	South BoundaryCondition `json:"south"`
	West  BoundaryCondition `json:"west"`
}

type BoundaryHandler interface {
	HandleBoundaryCoefficients(i, j, nx, ny int, diffX, diffY float64) (aE, aW, aN, aS float64)
	HandleBoundaryRHS(i, j, nx, ny int, aE, aW, aN, aS, dx, dy float64) float64
}

func NewDirichletBCs() *BoundaryConditions {
	return &BoundaryConditions{
		North: BoundaryCondition{Type: Dirichlet, Value: 0},
		East:  BoundaryCondition{Type: Dirichlet, Value: 0},
		South: BoundaryCondition{Type: Dirichlet, Value: 0},
		West:  BoundaryCondition{Type: Dirichlet, Value: 0},
	}
}

func (bc *BoundaryConditions) HandleBoundaryCoefficients(i, j, nx, ny int, diffX, diffY float64) (aE, aW, aN, aS float64) {
	aE = diffX
	aW = diffX
	aN = diffY
	aS = diffY

	if j == nx-1 { // East boundary
		aE = bc.getCoefficient(East, diffX)
	}
	if j == 0 { // West boundary
		aW = bc.getCoefficient(West, diffX)
	}
	if i == 0 { // North boundary
		aN = bc.getCoefficient(North, diffY)
	}
	if i == ny-1 { // South boundary
		aS = bc.getCoefficient(South, diffY)
	}

	return
}

func (bc *BoundaryConditions) HandleBoundaryRHS(i, j, nx, ny int, aE, aW, aN, aS, dx, dy float64) float64 {
	rhs := 0.0

	if j == nx-1 { // East boundary
		rhs += bc.getRHSContribution(East, aE, dy)
	}
	if j == 0 { // West boundary
		rhs += bc.getRHSContribution(West, aW, dy)
	}
	if i == 0 { // North boundary
		rhs += bc.getRHSContribution(North, aN, dx)
	}
	if i == ny-1 { // South boundary
		rhs += bc.getRHSContribution(South, aS, dx)
	}

	return rhs
}

func (bc *BoundaryConditions) SetNorth(bt BoundaryType, value float64 ) {
	bc.North = BoundaryCondition{ Type: bt, Value: value }
}

func (bc *BoundaryConditions) SetEast(bt BoundaryType, value float64 ) {
	bc.East = BoundaryCondition{ Type: bt, Value: value }
}

func (bc *BoundaryConditions) SetSouth(bt BoundaryType, value float64 ) {
	bc.South = BoundaryCondition{ Type: bt, Value: value }
}

func (bc *BoundaryConditions) SetWest(bt BoundaryType, value float64 ) {
	bc.West = BoundaryCondition{ Type: bt, Value: value }
}

func (bc *BoundaryConditions) getCoefficient(dir Direction, baseDiffCoeff float64) float64 {
	condition := bc.getCondition(dir)

	switch condition.Type {
	case Dirichlet:
		return 2 * baseDiffCoeff
	case Neumann:
		return 0
	case Outflow:
		return 0
	default:
		return baseDiffCoeff
	}
}

func (bc *BoundaryConditions) getRHSContribution(dir Direction, coeff, area float64) float64 {
	condition := bc.getCondition(dir)

	switch condition.Type {
	case Dirichlet:
		return coeff * condition.Value
	case Neumann:
		return area * condition.Value
	case Outflow:
		return 0
	default:
		return 0
	}
}

func (bc *BoundaryConditions) getCondition(dir Direction) BoundaryCondition {
	switch dir {
	case North:
		return bc.North
	case East:
		return bc.East
	case South:
		return bc.South
	case West:
		return bc.West
	default:
		return BoundaryCondition{Type: Neumann, Value: 0}
	}
}
