// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"embed"
	"encoding/json"
	"image/png"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/*
var assets embed.FS

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
