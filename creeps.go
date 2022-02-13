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
	Frame        int
	LastMoved    int
	Sprite       *SpriteSheet
}

// Update handles game logic for a Creep
func (c *Creep) Update(g *Game) error {
	if c.Health <= 0 {
		return errors.New("Creep died")
	}

	c.LastMoved = (c.LastMoved + 1) % 10
	if c.LastMoved != 0 {
		return nil
	}

	c.navigateWaypoints(g)
	return nil
}

func (c *Creep) navigateWaypoints(g *Game) {
	targetSquare := g.MapData[c.NextWaypoint]
	targertCoords := image.Pt(targetSquare.X*7+4, targetSquare.Y*7+4+5)
	if targertCoords.X > c.Coords.X {
		c.Coords.X++
	}
	if targertCoords.X < c.Coords.X {
		c.Coords.X--
	}
	if targertCoords.Y > c.Coords.Y {
		c.Coords.Y++
	}
	if targertCoords.Y < c.Coords.Y {
		c.Coords.Y--
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
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.Coords.X-3), float64(c.Coords.Y-3))
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
