//go:build wasm

package main

import (
	"syscall/js"
)

const (
	canvasWidthPx  = 800
	canvasHeightPx = 600
	canvasScale    = 0.002
	nx             = 33
	ny             = 21
)

func main() {
	done := make(chan struct{}, 0)

	bc := NewDirichletBCs()
	bc.SetSouth(Dirichlet, 0.8)
	bc.SetEast(Neumann, 0.1)
	bc.SetWest(Neumann, -0.5)
	simulation := NewSimulation(canvasWidthPx, canvasHeightPx, nx, ny, canvasScale, bc)

	renderer := NewRenderer("cfd-canvas", 800, 600, nx, ny)
	stats := NewStats()

	simulation.AddSource(nx/3, ny/2)

	go animationLoop(simulation, renderer, stats)

	<-done
}

func animationLoop(simulation *Simulation, renderer *Renderer, stats *Stats) {
	var animationFrame js.Func
	animationFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		dt := 0.016

		simulation.Step(dt)
		renderer.Draw(simulation)
		stats.Update(simulation)

		// Request next frame immediately (unthrottled!)
		js.Global().Call("requestAnimationFrame", animationFrame)
		return nil
	})

	// start the loop
	js.Global().Call("requestAnimationFrame", animationFrame)
}
