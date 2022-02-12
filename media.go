// Copyright 2022 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"embed"
	"encoding/json"
	"image/png"
	"io/ioutil"
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/*
var assets embed.FS

// NewMusicPlayer loads a sound into an audio player that can be used to play it
// as an infinite loop of music without any additional setup required
func NewMusicPlayer(music *vorbis.Stream, context *audio.Context) *audio.Player {
	musicLoop := audio.NewInfiniteLoop(music, music.Length())
	musicPlayer, err := audio.NewPlayer(context, musicLoop)
	if err != nil {
		log.Fatalf("error making music player: %v\n", err)
	}
	return musicPlayer
}

// NewSoundPlayer loads a sound into an audio player that can be used to play it
// without any additional setup required
func NewSoundPlayer(audioFile *vorbis.Stream, context *audio.Context) *audio.Player {
	audioPlayer, err := audio.NewPlayer(context, audioFile)
	if err != nil {
		log.Fatalf("error making audio player: %v\n", err)
	}
	return audioPlayer
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

// Frame is a single frame of an animation, usually a sub-image of a larger
// image containing several frames
type Frame struct {
	Duration int           `json:"duration"`
	Position FramePosition `json:"frame"`
}

// FramePosition represents the position of a frame, including the top-left
// coordinates and its dimensions (width and height)
type FramePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// Sprite is a slice of frames used to create sprite animation, it would've been
// nice to be able to call this "Frames" but that would cause confusion with the
// "frames" field in the JSON file
type Sprite []Frame

// SpriteSheet is the root-node of sprite data, it contains frames and meta data
// about them
type SpriteSheet struct {
	Sprite Sprite `json:"frames"`
	Image  *ebiten.Image
}

// Waypoint is a point marking a change of direction in the way along the map
type Waypoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Ways is a slice of waypoints from spawn point to the base
type Ways []*Waypoint

// MapData is waypoint data for a level map
type MapData struct {
	Ways Ways `json:"points"`
}

// Load map waypoint data from a given JSON file
func loadWays(name string) Ways {
	name = path.Join("assets", "maps", name)
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

	var mapdata MapData
	json.Unmarshal(data, &mapdata)
	if err != nil {
		log.Fatal(err)
	}

	return mapdata.Ways
}

// SoundType is a unique identifier to reference sound by name
type SoundType uint64

const (
	soundMusicTitle SoundType = iota
	soundMusicConstruction
	soundVictorious
	soundFail
)

// SpriteType is a unique identifier to load a sprite by name
type SpriteType uint64

const (
	spriteBigMonster SpriteType = iota
	spriteTowerBasic
	spriteTowerStrong
	spriteBigMonsterHorizont
	spriteBigMonsterVertical
	spriteBumm
	spriteSmallMonster
	spriteTinyMonster
	spriteTowerBottom
	spriteTowerLeft
	spriteTowerRight
	spriteTowerUp
	spriteHeartGone
	spriteIconHeart
	spriteIconMoney
	spriteIconTime
	spriteTitleScreen
)

// Load a sprite image and associated meta-data given a file name (without
// extension)
func loadSprite(name string) *SpriteSheet {
	name = path.Join("assets", "sprites", name)
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

	ss.Image = loadImage(name + ".png")

	return &ss
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

// Load a TTF font from a file in  embedded FS into a font face
func loadFont(name string, size float64) font.Face {
	log.Printf("loading %s\n", name)

	file, err := assets.Open(name)
	if err != nil {
		log.Fatalf("error opening file %s: %v\n", name, err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("error reading font file: ", err)
	}

	fontdata, err := opentype.Parse(data)
	if err != nil {
		log.Fatal("error parsing font data: ", err)
	}

	fontface, err := opentype.NewFace(fontdata, &opentype.FaceOptions{
		Size:    size, // The actual height of the font
		DPI:     72,   // This is a default, it looks horrible with any other value
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("error creating font face: ", err)
	}
	return fontface
}
