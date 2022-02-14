// Copyright 2022 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Cursor is used to interact with game entities at the given coordinates
type Cursor struct {
	Coords   image.Point
	Image    *ebiten.Image
	Cooldown int // Wait to show off construction animation
	Width    int
}

// Update implements Entity
func (c *Cursor) Update(g *Game) error {
	oldPos := c.Coords
	tileSize := 7
	hudOffset := 5

	if c.Cooldown > 0 {
		c.Cooldown--
	}

	// Movement controls
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		c.Move(image.Pt(0, tileSize))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		c.Move(image.Pt(0, -tileSize))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		c.Move(image.Pt(-tileSize, 0))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		c.Move(image.Pt(tileSize, 0))
	}

	// Keep the cursor inside the map
	if c.Coords.X < 0 ||
		c.Coords.Y < hudOffset ||
		c.Coords.X > g.Size.X ||
		c.Coords.Y > g.Size.Y {
		c.Coords = oldPos
	}

	return nil
}

// Move moves the player upwards
func (c *Cursor) Move(dest image.Point) {
	c.Coords = c.Coords.Add(dest)
	c.Cooldown = 0
}

// Draw implements Entity
func (c *Cursor) Draw(g *Game, screen *ebiten.Image) {
	if c.Cooldown != 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64(c.Coords.X-c.Width/2),
		float64(c.Coords.Y-c.Width/2),
	)
	screen.DrawImage(c.Image, op)
}

// NewCursor creates a new cursor struct at the bottom-left of the map
// It is shaped like a crosshair and is used to interact with the game
func NewCursor() *Cursor {
	tileSize := 7
	hudOffset := 6
	tileCenter := 3
	coords := image.Pt(
		2*tileSize+tileCenter,
		5*tileSize+tileCenter+hudOffset,
	)

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
