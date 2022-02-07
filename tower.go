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

// Towers is a slice of Tower entities
type Towers []*Tower
