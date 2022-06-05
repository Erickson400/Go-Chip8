package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// make an ebiten project simple.

type Game struct{}

var cpu CPU = CPU{}

func (g *Game) Update() error {
	cpu.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	cpu.Draw(screen)

	//rect := ebiten.NewImage(64, 32)
	//rect.Fill(color.RGBA{255, 255, 255, 255})
	//screen.DrawImage(rect, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	//size := 0.3
	//return int(1020 * size), int(700 * size)
	return 64, 32
}

func init() {
	err := cpu.Init("games/Rush Hour [Hap, 2006].ch8")
	//err := cpu.Init("test.ch8")
	if err != nil {
		panic(err)
	}
}

func main() {
	ebiten.SetWindowSize(1024, 512)
	ebiten.SetWindowTitle("Go Chip8")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}

}
