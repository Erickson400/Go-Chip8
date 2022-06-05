package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Display struct {
	//64x32
	buffer *ebiten.Image
}

func (d *Display) Init() {
	d.buffer = ebiten.NewImage(64, 32)
}

func (d *Display) Render(screen *ebiten.Image) {
	screen.DrawImage(d.buffer, nil)
}

func (d *Display) FlipPixel(x, y int, isTrue bool) bool {
	// leave if not want to change, or if pixel outisde of Image
	if !isTrue {
		return false
	}

	if x < 0 || x > 63 || y < 0 || y > 31 {
		return false
	}

	BLACK := color.RGBA{0, 0, 0, 0}
	WHITE := color.RGBA{255, 255, 255, 255}
	if d.buffer.At(x, y) == BLACK {
		d.buffer.Set(x, y, WHITE)
		return false
	} else {
		d.buffer.Set(x, y, BLACK)
		return true
	}
}

func (d *Display) Clear() {
	d.buffer.Clear()
}
