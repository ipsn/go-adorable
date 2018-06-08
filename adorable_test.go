// go-adorable - Adorable Avatars from Go
// Copyright (c) 2018 Péter Szilágyi. All rights reserved.

package adorable

import (
	"image/color"
	"testing"
)

// Tests just a bunch of random avatar generation.
func TestGenerate(t *testing.T) {
	for i := 0; i < 10; i++ {
		Random()
		RandomWithColor(color.RGBA{R: 64, G: 128, B: 192, A: 255})
		PseudoRandom([]byte("seed"))
		PseudoRandomWithColor([]byte("seed"), color.RGBA{R: 64, G: 128, B: 192, A: 255})
	}
}
