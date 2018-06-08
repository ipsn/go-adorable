// go-adorable - Adorable Avatars from Go
// Copyright (c) 2018 Péter Szilágyi. All rights reserved.

// Package adorable generates (pseudo) random user avatars.
package adorable

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"math/big"

	"github.com/ipsn/go-adorable/internal/bodyparts"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"
)

// Random creates a brand new random avatar and returns it in PNG format.
func Random() []byte {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	return PseudoRandom(seed)
}

// RandomWithColor creates a brand new random avatar using a predefined background
// color and returns it in PNG format.
func RandomWithColor(color color.RGBA) []byte {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	return PseudoRandomWithColor(seed, color)
}

// Random creates a brand new deterministic random avatar based on the provided
// seed and returns it in PNG format.
func PseudoRandom(seed []byte) []byte {
	random := hkdf.New(sha3.New256, seed, nil, nil)
	color := randomColor(random, 0.5, 1, 0.5, 0.9)

	return pseudoRandomWithColor(random, color)
}

// PseudoRandomWithColor creates a brand new deterministic random avatar based on
// the provided seed using a predefined background color and returns it in PNG
// format.
func PseudoRandomWithColor(seed []byte, color color.RGBA) []byte {
	random := hkdf.New(sha3.New256, seed, nil, nil)
	randomColor(random, 0.5, 1, 0.5, 0.9) // skip the color randomness

	return pseudoRandomWithColor(random, color)
}

// pseudoRandomWithColor creates a brand new deterministic random avatar based on
// the provided seed using a predefined background color and returns it in PNG
// format.
func pseudoRandomWithColor(random io.Reader, color color.RGBA) []byte {
	eye := randomImage(random, "eyes", 9)
	nose := randomImage(random, "nose", 8)
	mouth := randomImage(random, "mouth", 8)

	avatar := image.NewRGBA(eye.Bounds())

	draw.Draw(avatar, avatar.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	draw.Draw(avatar, avatar.Bounds(), eye, image.ZP, draw.Over)
	draw.Draw(avatar, avatar.Bounds(), nose, image.ZP, draw.Over)
	draw.Draw(avatar, avatar.Bounds(), mouth, image.ZP, draw.Over)

	blob := new(bytes.Buffer)
	if err := png.Encode(blob, avatar); err != nil {
		panic(err)
	}
	return blob.Bytes()
}

// randomFloat generates a single random floating point number based on the source
// randomness stream.
func randomFloat(random io.Reader) float64 {
	n, err := rand.Int(random, new(big.Int).SetUint64(math.MaxUint64))
	if err != nil {
		panic(err) // we can never exceed the HKDF maximum entropy
	}
	return float64(n.Uint64()) / float64(math.MaxUint64)
}

// randomColor generates a random background color. Internally it uses the HSV
// colorspace to permit clean and vivid random colors.
func randomColor(random io.Reader, minSaturation, maxSaturation, minValue, maxValue float64) color.RGBA {
	hcl := colorful.Hcl(
		360*randomFloat(random),
		minSaturation+(maxSaturation-minSaturation)*randomFloat(random),
		minValue+(maxValue-minValue)*randomFloat(random),
	).Clamped()

	return color.RGBA{
		R: byte(255 * hcl.R),
		G: byte(255 * hcl.G),
		B: byte(255 * hcl.B),
		A: 255,
	}
}

// randomImage loads an image of the given type, randomly from the available ones.
func randomImage(random io.Reader, kind string, have int) image.Image {
	n, err := rand.Int(random, big.NewInt(int64(have)))
	if err != nil {
		panic(err) // we can never exceed the HKDF maximum entropy
	}
	blob := bodyparts.MustAsset(fmt.Sprintf("%s%d.png", kind, n))

	img, err := png.Decode(bytes.NewReader(blob))
	if err != nil {
		panic(err)
	}
	return img
}
