// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Tower can be placed at a position to shoot Creeps
type Tower struct {
	Coords image.Point
	Cost   int
	Frame  int
	Sprite *SpriteSheet
}

// NewBasicTower is a convenience wrapper to make a basic-looking tower
func NewBasicTower(g *Game) *Tower {
	sprite, ok := g.Sprites[spriteTowerBasic]
	if !ok {
		log.Fatal("Failed to retrieve basic tower from game resource map")
	}
	return &Tower{g.Cursor.Coords, 200, 0, sprite}
}

// NewStrongTower is a convenience wrapper to make a strong-looking tower
func NewStrongTower(g *Game) *Tower {
	var sprite *SpriteSheet
	sprite, ok := g.Sprites[spriteTowerStrong]
	if !ok {
		log.Fatal("Failed to retrieve strong tower from game resource map")
	}
	return &Tower{g.Cursor.Coords, 500, 0, sprite}
}

// Update handles game logic for towers
func (t *Tower) Update(g *Game) {
	if t.Frame < len(t.Sprite.Sprite)-1 {
		t.Frame++
	}
}

// Draw draws the Tower to the screen
func (t *Tower) Draw(g *Game, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.Coords.X-1), float64(t.Coords.Y-1))
	s := t.Sprite
	frame := s.Sprite[t.Frame]
	screen.DrawImage(s.Image.SubImage(image.Rect(
		frame.Position.X,
		frame.Position.Y,
		frame.Position.X+frame.Position.W,
		frame.Position.Y+frame.Position.H,
	)).(*ebiten.Image), op)
}

// Towers is a slice of Tower entities
type Towers []Entity
