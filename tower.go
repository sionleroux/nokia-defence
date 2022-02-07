// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Tower can be placed at a position to shoot Creeps
type Tower struct {
	Coords image.Point
	Image  *ebiten.Image
}

// NewTower makes a new tower provided image pixel input
func NewTower(coords image.Point, size int, pix []uint8) *Tower {
	i := image.NewPaletted(
		image.Rect(0, 0, size, size),
		NokiaPalette,
	)
	i.Pix = pix

	return &Tower{coords, ebiten.NewImageFromImage(i)}
}

// NewBasicTower is a convenience wrapper to make a basic-looking tower
func NewBasicTower(coords image.Point) *Tower {
	return NewTower(coords, 5, []uint8{
		2, 2, 2, 2, 2,
		2, 1, 1, 1, 2,
		2, 1, 1, 1, 2,
		2, 1, 1, 1, 2,
		2, 2, 2, 2, 2,
	})
}

// NewStrongTower is a convenience wrapper to make a strong-looking tower
func NewStrongTower(coords image.Point) *Tower {
	return NewTower(coords, 5, []uint8{
		2, 2, 2, 2, 2,
		2, 2, 1, 2, 2,
		2, 1, 1, 1, 2,
		2, 2, 1, 2, 2,
		2, 2, 2, 2, 2,
	})
}

// Update handles game logic for towers
func (t *Tower) Update(g *Game) {
	panic("not implemented") // TODO: Implement
}

// Draw draws the Tower to the screen
func (t *Tower) Draw(g *Game, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.Coords.X-1), float64(t.Coords.Y-1))
	screen.DrawImage(t.Image, op)
}

// Towers is a slice of Tower entities
type Towers []Entity
