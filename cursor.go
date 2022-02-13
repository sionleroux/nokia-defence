// Copyright 2022 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Cursor is used to interact with game entities at the given coordinates
type Cursor struct {
	Coords image.Point
	Image  *ebiten.Image
	Width  int
}

// Update implements Entity
func (c *Cursor) Update(g *Game) error {
	// Movement controls
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		c.Move(image.Pt(0, 1))
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		c.Move(image.Pt(0, -1))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		c.Move(image.Pt(-1, 0))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		c.Move(image.Pt(1, 0))
	}

	return nil
}

// Move moves the player upwards
func (c *Cursor) Move(dest image.Point) {
	c.Coords = c.Coords.Add(dest)
}

// Draw implements Entity
func (c *Cursor) Draw(g *Game, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64(c.Coords.X-c.Width/2),
		float64(c.Coords.Y-c.Width/2),
	)
	screen.DrawImage(c.Image, op)
}

// NewCursor creates a new cursor struct at the given coordinates
// It is shaped like a crosshair and is used to interact with the game
func NewCursor(coords image.Point) *Cursor {

	w := 3
	i := image.NewPaletted(
		image.Rect(0, 0, w, w),
		NokiaPalette,
	)
	i.Pix = []uint8{
		0, 1, 0,
		1, 0, 1,
		0, 1, 0,
	}

	return &Cursor{
		Coords: coords,
		Image:  ebiten.NewImageFromImage(i),
		Width:  w,
	}
}
