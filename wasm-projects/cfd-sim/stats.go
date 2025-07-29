//go:build wasm

package main

import (
	"fmt"
	"syscall/js"
	"time"
)

type Stats struct {
	fpsElement    js.Value
	timeElement   js.Value
	maxPhiElement js.Value

	frameCount    int
	lastFPSUpdate time.Time
	startTime     time.Time
}

func NewStats() *Stats {
	doc := js.Global().Get("document")
	return &Stats{
		fpsElement:    doc.Call("getElementById", "cfd-fps"),
		timeElement:   doc.Call("getElementById", "cfd-time"),
		maxPhiElement: doc.Call("getElementById", "cfd-max-concentration"),
		lastFPSUpdate: time.Now(),
		startTime:     time.Now(),
	}
}

func (s *Stats) Update(simulation *Simulation) {
	s.frameCount++
	now := time.Now()

	// Update FPS every 500ms
	if now.Sub(s.lastFPSUpdate) >= 500*time.Millisecond {
		fps := float64(s.frameCount) / now.Sub(s.lastFPSUpdate).Seconds()
		s.fpsElement.Set("textContent", fmt.Sprintf("FPS: %.1f", fps))
		s.frameCount = 0
		s.lastFPSUpdate = now
	}

	// Update other stats every frame
	elapsedTime := now.Sub(s.startTime).Seconds()
	s.timeElement.Set("textContent", fmt.Sprintf("Time: %.2fs", elapsedTime))

	// Find max concentration
	maxPhi := 0.0
	for i := 0; i < simulation.mesh.ny; i++ {
		for j := 0; j < simulation.mesh.nx; j++ {
			if simulation.mesh.phi[i][j] > maxPhi {
				maxPhi = simulation.mesh.phi[i][j]
			}
		}
	}
	s.maxPhiElement.Set("textContent", fmt.Sprintf("Max φ: %.3f", maxPhi))
}
