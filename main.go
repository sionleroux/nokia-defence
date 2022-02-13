// Copyright 2022 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

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
	StartingMoney int = 500
)

func main() {
	windowScale := 10
	ebiten.SetWindowSize(GameSize.X*windowScale, GameSize.Y*windowScale)
	ebiten.SetWindowTitle("Nokia Defence")

	// Fonts
	font := loadFont("assets/fonts/tiny.ttf", 6)

	game := &Game{
		Loading: true,
		Size:    GameSize,
		Money:   StartingMoney,
		Font:    font,
	}

	go NewGame(game)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Game represents the main game state
type Game struct {
	Loading       bool
	Size          image.Point
	Cursor        *Cursor
	Maps          []*ebiten.Image
	MapData       Ways
	Sounds        []*vorbis.Stream
	Mplayer       []*vorbis.Stream
	Mcontext      *audio.Context
	Music         *audio.Player
	MapIndex      int
	Sprites       map[SpriteType]*SpriteSheet
	Towers        Towers
	Creeps        Creeps
	SpawnCooldown int
	Money         int
	Count         int
	Font          font.Face
}

// NewGame sets up a new game object with default states and game objects
func NewGame(g *Game) {

	// Music
	const sampleRate int = 44100 // assuming "normal" sample rate
	g.Mcontext = audio.NewContext(sampleRate)
	g.Sounds = make([]*vorbis.Stream, 4)
	g.Sounds[soundMusicConstruction] = loadSoundFile("assets/music/construction.ogg", sampleRate)
	g.Sounds[soundMusicTitle] = loadSoundFile("assets/music/title.ogg", sampleRate)
	g.Sounds[soundVictorious] = loadSoundFile("assets/sfx/victorious.ogg", sampleRate)
	g.Sounds[soundFail] = loadSoundFile("assets/sfx/fail.ogg", sampleRate)
	music := NewMusicPlayer(g.Sounds[soundMusicConstruction], g.Mcontext)
	music.Play()
	g.Music = music

	// Sprites
	g.Sprites = make(map[SpriteType]*SpriteSheet, 12)
	g.Sprites[spriteTowerBasic] = loadSprite("basic-tower")
	g.Sprites[spriteTowerStrong] = loadSprite("strong-tower")
	g.Sprites[spriteBigMonsterHorizont] = loadSprite("big_monster_horizont")
	g.Sprites[spriteBigMonsterVertical] = loadSprite("big_monster_vertical")
	g.Sprites[spriteSmallMonster] = loadSprite("small_monster")
	g.Sprites[spriteTinyMonster] = loadSprite("tiny_monster")
	g.Sprites[spriteBumm] = loadSprite("bumm")
	g.Sprites[spriteTowerBottom] = loadSprite("tower_bottom")
	g.Sprites[spriteTowerLeft] = loadSprite("tower_left")
	g.Sprites[spriteTowerRight] = loadSprite("tower_right")
	g.Sprites[spriteTowerUp] = loadSprite("tower_up")
	g.Sprites[spriteHeartGone] = loadSprite("heart_gone")
	g.Sprites[spriteIconHeart] = loadSprite("heart_icon")
	g.Sprites[spriteIconMoney] = loadSprite("money_icon")
	g.Sprites[spriteIconTime] = loadSprite("time_icon")
	g.Sprites[spriteTitleScreen] = loadSprite("titlescreen")

	// Static images
	g.Maps = make([]*ebiten.Image, 3)
	g.Maps[0] = loadImage("assets/maps/map1.png")
	g.Maps[1] = loadImage("assets/maps/map2.png")
	g.Maps[2] = loadImage("assets/maps/map3.png")
	g.MapData = loadWays("map1")

	g.Cursor = NewCursor(image.Pt(GameSize.X/2, GameSize.Y/2))

	g.Loading = false
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

	// Skip updating while the game is loading
	if g.Loading {
		return nil
	}

	g.Cursor.Update(g)

	for _, t := range g.Towers {
		t.Update(g)
	}

	for i, c := range g.Creeps {
		if err := c.Update(g); err != nil {
			log.Println(err)
			g.Creeps = append(g.Creeps[:i], g.Creeps[i+1:]...)
		}
	}

	// Tower placement controls
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		t := NewBasicTower(g)
		moneydiff := g.Money - t.Cost
		log.Printf("Buying tower %d - %d = %d\n", g.Money, t.Cost, moneydiff)
		if moneydiff >= 0 {
			g.Towers = append(g.Towers, t)
			g.Money = moneydiff
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		t := NewStrongTower(g)
		moneydiff := g.Money - t.Cost
		log.Printf("Buying tower %d - %d = %d\n", g.Money, t.Cost, moneydiff)
		if moneydiff >= 0 {
			g.Towers = append(g.Towers, t)
			g.Money = moneydiff
		}
	}

	if g.SpawnCooldown == 0 {
		spawn := g.MapData[0]
		gridScale := 7
		hudMargin := 5
		gridSquareMid := 4
		g.Creeps = append(g.Creeps, &Creep{
			Coords: image.Pt(
				spawn.X*gridScale+gridSquareMid,
				spawn.Y*gridScale+hudMargin+gridSquareMid,
			),
			NextWaypoint: 1,
			Damage:       0,
			Health:       1000,
			Loot:         50,
			Frame:        0,
			Sprite:       g.Sprites[spriteSmallMonster],
		})
	}

	// Spawn a new creep every N ticks
	g.SpawnCooldown = (g.SpawnCooldown + 1) % (3 * 60)

	return nil
}

// Draw draws the game screen by one frame
func (g *Game) Draw(screen *ebiten.Image) {
	// Light background
	screen.Fill(ColorLight)

	if g.Loading {
		// Try using text with pixel font
		txt := "Loading..."
		txtf, _ := font.BoundString(g.Font, txt)
		txth := (txtf.Max.Y - txtf.Min.Y).Ceil() / 2
		txtw := (txtf.Max.X - txtf.Min.X).Ceil() / 2
		text.Draw(screen, txt, g.Font, g.Size.X/2-txtw, g.Size.Y/2-txth, ColorDark)
		return
	}

	// Map background image
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.Maps[g.MapIndex], op)

	hudSize := 6.0
	ebitenutil.DrawRect(screen, 0, 0, float64(g.Size.X), hudSize, ColorDark)
	moneytxt := fmt.Sprintf("D%d", g.Money)
	text.Draw(screen, moneytxt, g.Font, 1, 5, ColorLight)

	for _, t := range g.Towers {
		t.Draw(g, screen)
	}

	for _, c := range g.Creeps {
		c.Draw(g, screen)
	}

	g.Cursor.Draw(g, screen)
}

// Entity is anything that can be interacted with in the game and drawn  to the
// screen, like Towers and Creeps
type Entity interface {
	Update(g *Game) error
	Draw(g *Game, screen *ebiten.Image)
}
