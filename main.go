// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"errors"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Media settings based on the Nokia 3310 jam restrictions
var (
	// ColorLight is the ON or 1 screen colour, similar to white
	ColorLight color.Color = color.RGBA{199, 240, 216, 255}
	// ColorDark is the OFF or 0 screen colour, similar to black
	ColorDark color.Color = color.RGBA{67, 82, 61, 255}
	// NokiaPalette is a 1-bit palette of greenish colours simulating Nokia 3310
	NokiaPalette color.Palette = color.Palette{ColorDark, ColorLight}
	// GameSize is the screen resolution of a Nokia 3310
	GameSize image.Point = image.Point{84, 48}
)

func main() {
	windowScale := 10
	ebiten.SetWindowSize(GameSize.X*windowScale, GameSize.Y*windowScale)
	ebiten.SetWindowTitle("Nokia Defence")

	game := &Game{
		Size:   GameSize,
		Cursor: NewCursor(image.Pt(GameSize.X/2, GameSize.Y/2)),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Game represents the main game state
type Game struct {
	Size   image.Point
	Cursor *Cursor
	Towers Towers
}

// Layout is hardcoded for now, may be made dynamic in future
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Size.X, g.Size.Y
}

// Update calculates game logic
func (g *Game) Update() error {

	// Pressing Q any time quits immediately
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return errors.New("game quit by player")
	}

	// Pressing F toggles full-screen
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		} else {
			ebiten.SetFullscreen(true)
		}
	}

	// Movement controls
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Cursor.Move(image.Pt(0, 1))
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Cursor.Move(image.Pt(0, -1))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Cursor.Move(image.Pt(-1, 0))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Cursor.Move(image.Pt(1, 0))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.Towers = append(g.Towers, &Tower{g.Cursor.Coords})
	}

	return nil
}

// Draw draws the game screen by one frame
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(ColorDark)
	for _, t := range g.Towers {
		drawSquare(screen, t.Coords, 2, ColorLight)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.Cursor.Coords.X), float64(g.Cursor.Coords.Y))
	screen.DrawImage(g.Cursor.Image, op)
}

// Cursor is used to interact with game entities at the given coordinates
type Cursor struct {
	Coords image.Point
	Image  *ebiten.Image
}

// NewCursor creates a new cursor struct at the given coordinates
// It is shaped like a crosshair and is used to interact with the game
func NewCursor(coords image.Point) *Cursor {

	i := image.NewPaletted(
		image.Rect(0, 0, 3, 3),
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
	}
}

// Move moves the player upwards
func (c *Cursor) Move(dest image.Point) {
	c.Coords = c.Coords.Add(dest)
}

// Tower can be placed at a position to shoot Creeps
type Tower struct {
	Coords image.Point
}

// Towers is a slice of Tower entities
type Towers []*Tower

// drawSquare is a convenience wrapper mapping point coordinates to float values
// and reducing repetition when the rectangles sides have a ratio of 1:1
func drawSquare(dest *ebiten.Image, pos image.Point, size int, color color.Color) {
	ebitenutil.DrawRect(
		dest,
		float64(pos.X),
		float64(pos.Y),
		float64(size),
		float64(size),
		color,
	)
}
