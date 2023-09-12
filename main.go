// Copyright 2022 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// Media settings based on the Nokia 3310 jam restrictions
var (
	// ColorTransparent is completely transparent, used for images that aren't
	// square shaped to show the underlying colour
	ColorTransparent color.Color = color.RGBA{0, 0, 0, 0}
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
		Size:  GameSize,
		Money: StartingMoney,
		Font:  font,
	}

	go NewGame(game)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Game represents the main game state
type Game struct {
	State         int
	Size          image.Point
	Cursor        *Cursor
	Maps          []*ebiten.Image
	MapData1      MapData
	MapData2      MapData
	Waves         []Creeps
	MapData       Ways
	NoBuild       NoBuild // Places where you can't build
	Sounds        []*audio.Player
	MapIndex      int
	Sprites       map[SpriteType]*SpriteSheet
	Towers        Towers
	Creeps        Creeps
	Spawned       int
	SpawnCooldown int
	Money         int
	Count         int
	TitleFrame    int
	Font          font.Face
}

const (
	gameStateLoading int = iota
	gameStateTitle
	gameStateBuild
	gameStateWon
	gameStateLose
	gameStateWin
	gameStateWaiting
	gameStatePause
)

// NewGame sets up a new game object with default states and game objects
func NewGame(g *Game) {

	// Music
	const sampleRate int = 44100 // assuming "normal" sample rate
	context := audio.NewContext(sampleRate)
	g.Sounds = make([]*audio.Player, 4)
	g.Sounds[soundMusicConstruction] = NewMusicPlayer(loadSoundFile("assets/music/construction.ogg", sampleRate), context)
	g.Sounds[soundMusicTitle] = NewMusicPlayer(loadSoundFile("assets/music/title.ogg", sampleRate), context)
	g.Sounds[soundVictorious] = NewSoundPlayer(loadSoundFile("assets/sfx/victorious.ogg", sampleRate), context)
	g.Sounds[soundFail] = NewSoundPlayer(loadSoundFile("assets/sfx/fail.ogg", sampleRate), context)
	g.Sounds[soundMusicTitle].Play()

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
	g.MapData1 = loadWays("map1")
	g.MapData2 = loadWays("map2")
	g.MapData = g.MapData1.Ways
	g.NoBuild = g.MapData1.NoBuild

	g.Waves = NewWaves(g)
	g.Cursor = NewCursor()

	g.State = gameStateTitle
}

// Reset the game to initial state, ready for a new round
func (g *Game) Reset(win bool) {
	g.Creeps = nil
	g.Towers = nil
	g.SpawnCooldown = 0
	g.Spawned = 0
	g.Waves = NewWaves(g)
	g.Money = StartingMoney
	g.Count = 0
	g.TitleFrame = 0
	g.Cursor = NewCursor()
	if win && g.MapIndex < 1 {
		g.State = gameStateWaiting
		g.MapData = g.MapData2.Ways
		g.NoBuild = g.MapData2.NoBuild
		g.MapIndex++
		g.Sounds[soundMusicConstruction].Play()
		g.State = gameStateBuild
	} else {
		g.MapData = g.MapData1.Ways
		g.NoBuild = g.MapData1.NoBuild
		g.MapIndex = 0
		g.Sounds[soundMusicTitle].Play()
		if win {
			g.State = gameStateWon
		} else {
			g.State = gameStateTitle
		}
	}
}

// Layout is hardcoded for now, may be made dynamic in future
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Size.X, g.Size.Y
}

// Update calculates game logic
func (g *Game) Update() error {

	// Pressing F toggles full-screen
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		} else {
			ebiten.SetFullscreen(true)
		}
	}

	// Skip updating while the game is loading
	if g.State == gameStateLoading || g.State == gameStateWaiting {
		return nil
	}

	if g.State == gameStateWon && inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.State = gameStateTitle
		return nil
	}

	if g.State == gameStateLose {
		g.Sounds[soundMusicConstruction].Pause()
		g.Sounds[soundFail].Rewind()
		g.Sounds[soundFail].Play()
		g.State = gameStateWaiting
		gloat := time.NewTimer(time.Second * 4)
		go func() {
			log.Println("Gloating")
			<-gloat.C
			g.Reset(false)
		}()
		return nil
	}

	if g.State == gameStateWin {
		g.Sounds[soundMusicConstruction].Pause()
		g.Sounds[soundVictorious].Rewind()
		g.Sounds[soundVictorious].Play()
		g.State = gameStateWaiting
		gloat := time.NewTimer(time.Second * 2)
		go func() {
			log.Println("Gloating")
			<-gloat.C
			g.Reset(true)
		}()
		return nil
	}

	if g.State == gameStateTitle {
		g.Count = (g.Count + 1) % 15
		if g.Count == 0 {
			g.TitleFrame++
		}
		if g.TitleFrame > 19 {
			g.TitleFrame = 16 // XXX copied these from the JSON file cos I'm tired
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyX) {
			g.State = gameStateBuild
			g.Sounds[soundMusicTitle].Pause()
			g.Sounds[soundMusicConstruction].Play()
		}
		return nil
	}

	if g.State == gameStatePause {
		if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
			g.State = gameStateBuild
		}
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		g.State = gameStatePause
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

	if g.Spawned == len(g.Waves[g.MapIndex]) && len(g.Creeps) <= 0 {
		log.Println("You win")
		g.State = gameStateWin
	}

	// Tower placement controls
	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		BuyTower(g)
	}
	// Sell a tower
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if k := IsOccupied(g, g.Cursor.Coords); k != -1 {
			g.Towers = append(g.Towers[:k], g.Towers[k+1:]...)
			g.Money += 100
		}
	}

	if g.SpawnCooldown == 0 {
		spawn := g.MapData[0]
		gridScale := 7
		hudMargin := 5
		gridSquareMid := 4
		if g.Spawned < len(g.Waves[g.MapIndex]) {
			creep := g.Waves[g.MapIndex][g.Spawned]
			creep.Coords = image.Pt(
				spawn.X*gridScale+gridSquareMid,
				spawn.Y*gridScale+hudMargin+gridSquareMid,
			)
			g.Creeps = append(g.Creeps, creep)
			g.Spawned++
		}
	}

	// Spawn a new creep every N ticks
	g.SpawnCooldown = (g.SpawnCooldown + 1) % (3 * 60)

	return nil
}

// Draw draws the game screen by one frame
func (g *Game) Draw(screen *ebiten.Image) {
	// Light background
	screen.Fill(ColorLight)

	if g.State == gameStateLoading {
		txt := "Loading..."
		txtf, _ := font.BoundString(g.Font, txt)
		txth := (txtf.Max.Y - txtf.Min.Y).Ceil() / 2
		txtw := (txtf.Max.X - txtf.Min.X).Ceil() / 2
		text.Draw(screen, txt, g.Font, g.Size.X/2-txtw, g.Size.Y/2-txth, ColorDark)
		return
	}

	if g.State == gameStateWon {
		txt := "YOU WON!"
		txtf, _ := font.BoundString(g.Font, txt)
		txth := (txtf.Max.Y - txtf.Min.Y).Ceil() / 2
		txtw := (txtf.Max.X - txtf.Min.X).Ceil() / 2
		text.Draw(screen, txt, g.Font, g.Size.X/2-txtw, g.Size.Y/2-txth, ColorDark)
		return
	}

	if g.State == gameStatePause {
		txt := "Paused..."
		txtf, _ := font.BoundString(g.Font, txt)
		txth := (txtf.Max.Y - txtf.Min.Y).Ceil() / 2
		txtw := (txtf.Max.X - txtf.Min.X).Ceil() / 2
		text.Draw(screen, txt, g.Font, g.Size.X/2-txtw, g.Size.Y/2-txth, ColorDark)
		return
	}

	if g.State == gameStateTitle {
		s := g.Sprites[spriteTitleScreen]
		frame := s.Sprite[g.TitleFrame]
		screen.DrawImage(s.Image.SubImage(image.Rect(
			frame.Position.X,
			frame.Position.Y,
			frame.Position.X+frame.Position.W,
			frame.Position.Y+frame.Position.H,
		)).(*ebiten.Image), &ebiten.DrawImageOptions{})
		return
	}

	// Map background image
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.Maps[g.MapIndex], op)

	hudSize := 6.0
	ebitenutil.DrawRect(screen, 0, 0, float64(g.Size.X), hudSize, ColorDark)
	moneytxt := fmt.Sprintf("D%d", g.Money)
	text.Draw(screen, moneytxt, g.Font, 1, 5, ColorLight)
	var cost int
	if IsOccupied(g, g.Cursor.Coords) != -1 {
		cost = 300
	} else {
		cost = 200
	}
	costtxt := fmt.Sprintf("c%d", cost)
	costtxtf, _ := font.BoundString(g.Font, costtxt)
	costtxtw := (costtxtf.Max.X - costtxtf.Min.X).Ceil()
	text.Draw(screen, costtxt, g.Font, g.Size.X-costtxtw-1, 5, ColorLight)

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
