// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ImageBasicTower  *ebiten.Image
	ImageStrongTower *ebiten.Image
)

func init() {

	i := image.NewPaletted(
		image.Rect(0, 0, 5, 5),
		NokiaPalette,
	)

	i.Pix = []uint8{
		2, 2, 2, 2, 2,
		2, 1, 1, 1, 2,
		2, 1, 1, 1, 2,
		2, 1, 1, 1, 2,
		2, 2, 2, 2, 2,
	}
	ImageBasicTower = ebiten.NewImageFromImage(i)

	i.Pix = []uint8{
		2, 2, 2, 2, 2,
		2, 2, 1, 2, 2,
		2, 1, 1, 1, 2,
		2, 2, 1, 2, 2,
		2, 2, 2, 2, 2,
	}
	ImageStrongTower = ebiten.NewImageFromImage(i)

}

// Tower can be placed at a position to shoot Creeps
type Tower struct {
	Coords image.Point
}

// Update handles game logic for towers
func (t *Tower) Update(g *Game) {
	panic("not implemented") // TODO: Implement
}

// Draw draws the Tower to the screen
func (t *Tower) Draw(g *Game, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.Coords.X-1), float64(t.Coords.Y-1))
	screen.DrawImage(ImageBasicTower, op)
}

// Towers is a slice of Tower entities
type Towers []Entity
