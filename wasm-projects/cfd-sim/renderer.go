//go:build wasm

package main

import (
	"fmt"
	"syscall/js"
)

type Renderer struct {
	canvas      js.Value
	ctx         js.Value
	width       int
	height      int
	cellWidth   float64
	cellHeight  float64
	smoothedMax float64
	alpha       float64
}

func NewRenderer(canvasID string, canvasWidth, canvasHeight, nx, ny int) *Renderer {
	canvas := js.Global().Get("document").Call("getElementById", canvasID)
	ctx := canvas.Call("getContext", "2d")

	return &Renderer{
		canvas:     canvas,
		ctx:        ctx,
		width:      canvasWidth,
		height:     canvasHeight,
		cellWidth:  float64(canvasWidth) / float64(nx),
		cellHeight: float64(canvasHeight) / float64(ny),
		smoothedMax: 1.0,
		alpha:      0.01,
	}
}

func (r *Renderer) Draw(simulation *Simulation) {
	// Clear canvas
	r.ctx.Set("fillStyle", "white")
	r.ctx.Call("fillRect", 0, 0, r.width, r.height)

	// Find current max
	currentMax := 0.0
	for i := range simulation.mesh.ny {
		for j := range simulation.mesh.nx {
			if simulation.mesh.phi[i][j] > currentMax {
				currentMax = simulation.mesh.phi[i][j]
			}
		}
	}

	// exponential smoothing
	r.smoothedMax = r.alpha*currentMax + (1-r.alpha)*r.smoothedMax

	// Get Math.ceil function
	mathCeil := js.Global().Get("Math").Get("ceil")

	// Draw each cell using Math.ceil for dimensions
	for i := range simulation.mesh.ny {
		for j := range simulation.mesh.nx {
			phi := simulation.mesh.phi[i][j]
			normalisedPhi := phi / r.smoothedMax

			greyValue := int(255 * (1.0 - normalisedPhi))
			color := fmt.Sprintf("rgb(%d,%d,%d)", greyValue, greyValue, greyValue)
			r.ctx.Set("fillStyle", color)

			x := float64(j) * r.cellWidth
			y := float64(i) * r.cellHeight

			// Use Math.ceil to ensure we cover any fractional pixels
			width := mathCeil.Invoke(r.cellWidth).Float()
			height := mathCeil.Invoke(r.cellHeight).Float()

			r.ctx.Call("fillRect", x, y, width+1, height+1)
		}
	}
}
