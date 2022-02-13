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
	Damage int
	Frame  int
	Target *Creep // the creep it's currently attacking
	Sprite *SpriteSheet
}

// NewBasicTower is a convenience wrapper to make a basic-looking tower
func NewBasicTower(g *Game) *Tower {
	sprite, ok := g.Sprites[spriteTowerBasic]
	if !ok {
		log.Fatal("Failed to retrieve basic tower from game resource map")
	}
	return &Tower{g.Cursor.Coords, 200, 2, 0, nil, sprite}
}

// NewStrongTower is a convenience wrapper to make a strong-looking tower
func NewStrongTower(g *Game) *Tower {
	var sprite *SpriteSheet
	sprite, ok := g.Sprites[spriteTowerStrong]
	if !ok {
		log.Fatal("Failed to retrieve strong tower from game resource map")
	}
	return &Tower{g.Cursor.Coords, 500, 5, 0, nil, sprite}
}

// Update handles game logic for towers
func (t *Tower) Update(g *Game) error {
	// Construction animation
	if t.Frame < len(t.Sprite.Sprite)-1 {
		t.Frame++
	}

	// Target Seeking
	if t.Target == nil {
		t.findNewTarget(g)
	} else {
		t.clearIfOutOfRange()
	}

	// Damage dealing
	if t.Target != nil {
		died := t.Target.Attack(t.Damage)
		if died {
			t.Target = nil
		}
	}

	return nil
}

// Look for the first creep in range
func (t *Tower) findNewTarget(g *Game) {
	tileSize := 7
	rangeSize := 2 * tileSize
	for _, v := range g.Creeps {
		hitboxRadius := 3
		creepBox := image.Rectangle{
			v.Coords.Add(image.Pt(-hitboxRadius, -hitboxRadius)),
			v.Coords.Add(image.Pt(hitboxRadius, hitboxRadius)),
		}
		towerBox := image.Rect(
			t.Coords.X-rangeSize,
			t.Coords.Y-rangeSize,
			t.Coords.X+rangeSize,
			t.Coords.Y+rangeSize,
		)
		withinRange := towerBox.Overlaps(creepBox)
		if withinRange {
			t.Target = v
		}
	}
}

// Stop targeting a creep if it's already dead
func (t *Tower) cullDeadCreep() {
	if t.Target.Health <= 0 {
		t.Target = nil
	}
}

// Clear current target when it gets out of range
func (t *Tower) clearIfOutOfRange() {
	tileSize := 7
	rangeSize := 2 * tileSize
	hitboxRadius := 3
	creepBox := image.Rectangle{
		t.Target.Coords.Add(image.Pt(-hitboxRadius, -hitboxRadius)),
		t.Target.Coords.Add(image.Pt(hitboxRadius, hitboxRadius)),
	}
	towerBox := image.Rect(
		t.Coords.X-rangeSize,
		t.Coords.Y-rangeSize,
		t.Coords.X+rangeSize,
		t.Coords.Y+rangeSize,
	)
	withinRange := towerBox.Overlaps(creepBox)
	if !withinRange {
		t.Target = nil
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
	if t.Target != nil {
		c := t.Target
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
type Towers []*Tower
