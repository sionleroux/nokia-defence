// Copyright 2022 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Tower can be placed at a position to shoot Creeps
type Tower struct {
	Coords image.Point
	Cost   int
	Frame  int
	Target int // the creep it's currently attacking
	Sprite *SpriteSheet
}

// NewBasicTower is a convenience wrapper to make a basic-looking tower
func NewBasicTower(g *Game) *Tower {
	sprite, ok := g.Sprites[spriteTowerBasic]
	if !ok {
		log.Fatal("Failed to retrieve basic tower from game resource map")
	}
	return &Tower{g.Cursor.Coords, 200, 0, -1, sprite}
}

// NewStrongTower is a convenience wrapper to make a strong-looking tower
func NewStrongTower(g *Game) *Tower {
	var sprite *SpriteSheet
	sprite, ok := g.Sprites[spriteTowerStrong]
	if !ok {
		log.Fatal("Failed to retrieve strong tower from game resource map")
	}
	return &Tower{g.Cursor.Coords, 500, 0, -1, sprite}
}

// Update handles game logic for towers
func (t *Tower) Update(g *Game) {
	// Construction animation
	if t.Frame < len(t.Sprite.Sprite)-1 {
		t.Frame++
	}

	// Look for the first creep in range
	tileSize := 7
	for k, v := range g.Creeps {
		withinRange := image.Rectangle{
			v.(*Creep).Coords,
			v.(*Creep).Coords,
		}.In(
			image.Rect(
				t.Coords.X-tileSize,
				t.Coords.Y-tileSize,
				t.Coords.X+tileSize,
				t.Coords.Y+tileSize,
			))
		if withinRange {
			t.Target = k
		}
	}
}

// Draw draws the Tower to the screen
func (t *Tower) Draw(g *Game, screen *ebiten.Image) {

	// Draw tower
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

	// Draw shooting laser
	if t.Target != -1 {
		c := g.Creeps[t.Target].(*Creep)
		ebitenutil.DrawLine(screen,
			float64(t.Coords.X),
			float64(t.Coords.Y),
			float64(c.Coords.X),
			float64(c.Coords.Y),
			ColorDark,
		)
	}
}

// Towers is a slice of Tower entities
type Towers []Entity
