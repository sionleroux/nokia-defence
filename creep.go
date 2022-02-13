// Copyright 2022 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"errors"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Creep moves along a path from a spawn point towards the base it is attacking
type Creep struct {
	Coords       image.Point
	NextWaypoint int
	Health       int // Hit points
	Damage       int // How much damage it deals to the base
	Loot         int // How much money you get when it dies
	Frame        int
	LastMoved    int
	Direction    int  // Which way the creep is moving
	Flip         bool // Whether to flip the animation frame
	Sprite       *SpriteSheet
}

const (
	directionRight int = iota
	directionLeft
	directionUp
	directionDown
)

// Update handles game logic for a Creep
func (c *Creep) Update(g *Game) error {
	if c.Health <= 0 {
		g.Money += c.Loot
		return errors.New("Creep died")
	}

	c.LastMoved = (c.LastMoved + 1) % 10
	if c.LastMoved != 0 {
		return nil
	}

	c.navigateWaypoints(g)
	c.animate()

	return nil
}

func (c *Creep) animate() {
	const (
		HORIZONTAL = 0
		VERTICAL   = 2
	)
	var frameTag int
	switch c.Direction {
	case directionRight:
		c.Flip = false
		frameTag = HORIZONTAL
	case directionLeft:
		c.Flip = true
		frameTag = HORIZONTAL
	default:
		c.Flip = false
		frameTag = VERTICAL
	}
	from := c.Sprite.Meta.FrameTags[frameTag].From
	to := c.Sprite.Meta.FrameTags[frameTag].To
	if c.Frame < from || c.Frame >= to {
		c.Frame = from
		return
	}
	if c.Frame < to {
		c.Frame++
	}
}

func (c *Creep) navigateWaypoints(g *Game) {
	targetSquare := g.MapData[c.NextWaypoint]
	targertCoords := image.Pt(targetSquare.X*7+4, targetSquare.Y*7+4+5)
	if targertCoords.X > c.Coords.X {
		c.Coords.X++
		c.Direction = directionRight
	}
	if targertCoords.X < c.Coords.X {
		c.Coords.X--
		c.Direction = directionLeft
	}
	if targertCoords.Y > c.Coords.Y {
		c.Coords.Y++
		c.Direction = directionUp
	}
	if targertCoords.Y < c.Coords.Y {
		c.Coords.Y--
		c.Direction = directionDown
	}
	if targertCoords.X == c.Coords.X && targertCoords.Y == c.Coords.Y {
		next := c.NextWaypoint + 1
		if next < len(g.MapData) {
			c.NextWaypoint++
		} else {
			log.Fatal("You failed")
		}
	}
}

// Attack hurts a creep's health by a specified amount
func (c *Creep) Attack(amount int) bool {
	c.Health = c.Health - amount
	if c.Health <= 0 {
		return true
	}
	return false
}

// Draw draws the Creep to the screen
func (c *Creep) Draw(g *Game, screen *ebiten.Image) {
	s := c.Sprite
	frame := s.Sprite[c.Frame]
	op := &ebiten.DrawImageOptions{}
	if c.Flip { // Please don't ask
		op.GeoM.Translate(float64(-1*frame.Position.W/2), 1)
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(frame.Position.W/2), 1)
	}
	op.GeoM.Translate(float64(c.Coords.X-3), float64(c.Coords.Y-3))
	screen.DrawImage(s.Image.SubImage(image.Rect(
		frame.Position.X,
		frame.Position.Y,
		frame.Position.X+frame.Position.W,
		frame.Position.Y+frame.Position.H,
	)).(*ebiten.Image), op)
}

// Creeps is a slice of Creep entities
type Creeps []*Creep