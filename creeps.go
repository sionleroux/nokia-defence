// Copyright 2022 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Creep moves along a path from a spawn point towards the base it is attacking
type Creep struct {
	Coords image.Point
	Damage int
	Frame  int
	Sprite *SpriteSheet
}

// Update handles game logic for a Creep
func (c *Creep) Update(g *Game) {
	panic("not implemented") // TODO: Implement
}

// Draw draws the Creep to the screen
func (c *Creep) Draw(g *Game, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.Coords.X-1), float64(c.Coords.Y-1))
	s := c.Sprite
	frame := s.Sprite[c.Frame]
	screen.DrawImage(s.Image.SubImage(image.Rect(
		frame.Position.X,
		frame.Position.Y,
		frame.Position.X+frame.Position.W,
		frame.Position.Y+frame.Position.H,
	)).(*ebiten.Image), op)
}

// Creeps is a slice of Tower entities
type Creeps []Entity
