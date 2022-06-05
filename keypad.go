package main

import "github.com/hajimehoshi/ebiten/v2"

type Keypad struct {
	Keymap  [16]ebiten.Key // for changing the keys on your keyboard
	Keydata [16]bool       // to check if the key is pressed
}

func (k *Keypad) Init() {
	k.Keymap = [16]ebiten.Key{
		ebiten.KeyX,
		ebiten.Key1,
		ebiten.Key2,
		ebiten.Key3,
		ebiten.KeyQ,
		ebiten.KeyW,
		ebiten.KeyE,
		ebiten.KeyA,
		ebiten.KeyS,
		ebiten.KeyD,
		ebiten.KeyZ,
		ebiten.KeyC,
		ebiten.Key4,
		ebiten.KeyR,
		ebiten.KeyF,
		ebiten.KeyV,
	}

	for i := 0; i < 16; i++ {
		k.Keydata[i] = false
	}
}

func (k *Keypad) Update() {
	for i := 0; i < 16; i++ {
		k.Keydata[i] = ebiten.IsKeyPressed(k.Keymap[i])
	}
}
