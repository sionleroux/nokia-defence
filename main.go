// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"embed"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/*
var assets embed.FS

// Media settings based on the Nokia 3310 jam restrictions
var (
	// ColorTransparent is completely transparent, used for images that aren't
	// square shaped to show the underlying colour
	ColorTransparent color.Color = color.RGBA{67, 82, 61, 0}
	// ColorLight is the ON or 1 screen colour, similar to white
	ColorLight color.Color = color.RGBA{199, 240, 216, 255}
	// ColorDark is the OFF or 0 screen colour, similar to black
	ColorDark color.Color = color.RGBA{67, 82, 61, 255}
	// NokiaPalette is a 1-bit palette of greenish colours simulating Nokia 3310
	NokiaPalette color.Palette = color.Palette{ColorTransparent, ColorDark, ColorLight}
	// GameSize is the screen resolution of a Nokia 3310
	GameSize image.Point = image.Point{84, 48}
	// StartingMoney is the amount of money you start the game with
	StartingMoney int = 1000
)

func main() {
	windowScale := 10
	ebiten.SetWindowSize(GameSize.X*windowScale, GameSize.Y*windowScale)
	ebiten.SetWindowTitle("Nokia Defence")

	// Music
	const sampleRate int = 44100 // assuming "normal" sample rate
	music := loadSoundFile("assets/construction.ogg", sampleRate)
	musicLoop := audio.NewInfiniteLoop(music, music.Length())
	musicPlayer, err := audio.NewPlayer(audio.NewContext(sampleRate), musicLoop)
	if err != nil {
		log.Fatalf("error making music player: %v\n", err)
	}
	musicPlayer.Play()

	// Sprites
	basicsprite := loadSprite("assets/basic-tower")
	basicimage := loadImage("assets/basic-tower.png")

	game := &Game{
		Size:        GameSize,
		Cursor:      NewCursor(image.Pt(GameSize.X/2, GameSize.Y/2)),
		Money:       StartingMoney,
		BasicSprite: basicsprite,
		BasicImage:  basicimage,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Game represents the main game state
type Game struct {
	Size        image.Point
	Cursor      *Cursor
	Towers      Towers
	Money       int
	BasicSprite Sprite
	BasicImage  *ebiten.Image
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

	for _, t := range g.Towers {
		t.Update(g)
	}

	// Tower placement controls
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		t := NewBasicTower(g.Cursor.Coords)
		moneydiff := g.Money - t.Cost
		log.Printf("Buying tower %d - %d = %d\n", g.Money, t.Cost, moneydiff)
		if moneydiff >= 0 {
			g.Towers = append(g.Towers, t)
			g.Money = moneydiff
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		t := NewStrongTower(g.Cursor.Coords)
		moneydiff := g.Money - t.Cost
		log.Printf("Buying tower %d - %d = %d\n", g.Money, t.Cost, moneydiff)
		if moneydiff >= 0 {
			g.Towers = append(g.Towers, t)
			g.Money = moneydiff
		}
	}

	return nil
}

// Draw draws the game screen by one frame
func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.Fill(ColorLight)
	for _, t := range g.Towers {
		t.Draw(g, screen)
	}
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
		0, 2, 0,
		2, 0, 2,
		0, 2, 0,
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

// Entity is anything that can be interacted with in the game and drawn  to the
// screen, like Towers and Creeps
type Entity interface {
	Update(g *Game)
	Draw(g *Game, screen *ebiten.Image)
}

// Load an OGG Vorbis sound file with 44100 sample rate and return its stream
func loadSoundFile(name string, sampleRate int) *vorbis.Stream {
	log.Printf("loading %s\n", name)

	file, err := assets.Open(name)
	if err != nil {
		log.Fatalf("error opening file %s: %v\n", name, err)
	}
	defer file.Close()

	music, err := vorbis.DecodeWithSampleRate(sampleRate, file)
	if err != nil {
		log.Fatalf("error decoding file %s as Vorbis: %v\n", name, err)
	}

	return music
}

type Frame struct {
	Duration int           `json:"duration"`
	Position FramePosition `json:"frame"`
}

type FramePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type Sprite []Frame

type SpriteSheet struct {
	Sprite Sprite `json:"frames"`
}

// Load an OGG Vorbis sound file with 44100 sample rate and return its stream
func loadSprite(name string) Sprite {
	log.Printf("loading %s\n", name)

	file, err := assets.Open(name + ".json")
	if err != nil {
		log.Fatalf("error opening file %s: %v\n", name, err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var ss SpriteSheet
	json.Unmarshal(data, &ss)
	if err != nil {
		log.Fatal(err)
	}

	return ss.Sprite
}

// Load an image from embedded FS into an ebiten Image object
func loadImage(name string) *ebiten.Image {
	log.Printf("loading %s\n", name)

	file, err := assets.Open(name)
	if err != nil {
		log.Fatalf("error opening file %s: %v\n", name, err)
	}
	defer file.Close()

	raw, err := png.Decode(file)
	if err != nil {
		log.Fatalf("error decoding file %s as PNG: %v\n", name, err)
	}

	return ebiten.NewImageFromImage(raw)
}
