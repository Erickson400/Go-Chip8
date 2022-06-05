package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CYCLE_SIZE = 10 //10
	STACK_SIZE = 20
)

var DebugLogs = true
var LogFile *os.File

func init() {
	f, err := os.OpenFile("Logs.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	LogFile = f
	if err != nil {
		panic(err)
	}
}

func DebugLog(s string) {
	if DebugLogs {
		fmt.Fprintln(LogFile, s)
	}
}

//chip8 CPU
type CPU struct {
	Screen     *Display
	Pad        *Keypad
	Ram        [4095]byte
	Stack      [STACK_SIZE]uint16
	StackIndex int
	Register   [16]byte
	Delay      byte
	DelaySound byte
	IReg       uint16
	Pc         uint16
	IsRunning  bool
}

// TODO make the cpu functions after making display and keypad structs

func (c *CPU) Init(filename string) error {
	c.Screen = &Display{}
	c.Screen.Init()
	c.Pad = &Keypad{}
	c.Pad.Init()

	c.IsRunning = true
	c.Pc = 0x200

	// Read the file into the memory
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	font := [80]byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	//load the font into memory
	for i := 0; i < 80; i++ {
		c.Ram[i] = font[i]
	}

	//load the game into memory
	for i := 0; i < len(fileData); i++ {
		c.Ram[i+0x200] = fileData[i]
	}

	return nil
}

func (c *CPU) Update() {
	if !c.IsRunning {
		return
	}
	c.Pad.Update()

	for i := 0; i < CYCLE_SIZE; i++ {
		// Clocl Cycle
		c.Decode(c.Ram[c.Pc], c.Ram[c.Pc+1])
		if c.Delay > 0 {
			c.Delay--
		}
		if c.DelaySound > 0 {
			c.DelaySound--
		}
	}
}

func (c *CPU) Draw(screen *ebiten.Image) {
	if c.IsRunning {
		c.Screen.Render(screen)
	}
}

func (c *CPU) Decode(firstHalf, secondHalf byte) {

	nibble := [4]byte{ // nibble is a 4 bit value
		(firstHalf >> 4) & 0xF,
		firstHalf & 0xF,
		(secondHalf >> 4) & 0xF,
		secondHalf & 0xF,
	}

	// increment the program counter if less than 0xFFD
	if c.Pc < 0xFFD {
		c.Pc += 2
	} else {
		fmt.Println("WARNING: PC is at 0xFFD")
	}

	//Pre-calculated variables
	var VX byte = nibble[1]
	var VY byte = nibble[2]
	var N byte = nibble[3]
	var NN byte = (nibble[2] << 4) + nibble[3]
	var NNN uint16 = (uint16(nibble[1]) << 8) + uint16((nibble[2])<<4) + uint16(nibble[3])

	// execute the opcode
	switch nibble[0] {
	case 0x0:
		switch NNN {
		case 0xE0: // Clear screen
			c.Screen.Clear()
			DebugLog("Clear screen")
		case 0xEE: // Pop Stack
			if c.StackIndex > 0 {
				c.StackIndex--
				c.Pc = c.Stack[c.StackIndex]
				DebugLog("Pop stack")
			} else {
				fmt.Println("WARNING: Stack underflow")
			}
		default:
			fmt.Printf("WARNING: Unknown opcode: 0x%X\n", NNN)
		}
	case 0x1: // Set PC
		c.Pc = NNN
		DebugLog("Set PC " + fmt.Sprintf("(%X)", NNN))
	case 0x2: // Push to Stack
		if c.StackIndex < STACK_SIZE {
			c.Stack[c.StackIndex] = c.Pc
			c.StackIndex++
			c.Pc = NNN
			DebugLog("Push to stack")
		} else {
			fmt.Println("WARNING: Stack overflow")
		}
	case 0x3: // Skip if VX == NN
		if c.Register[VX] == NN {
			c.Pc += 2
			DebugLog("Skip if VX == NN")
		}
	case 0x4: // Skip if VX != NN
		if c.Register[VX] != NN {
			c.Pc += 2
			DebugLog("Skip if VX != NN")
		}
	case 0x5: // Skip if VX == VY
		if c.Register[VX] == c.Register[VY] {
			c.Pc += 2
			DebugLog("Skip if VX == VY")
		}
	case 0x6: // Set VX to NN
		c.Register[VX] = NN
		DebugLog("Set VX to NN " + fmt.Sprintf("(%X, %X)", VX, NN))
	case 0x7: // Add NN to VX
		c.Register[VX] += NN
		DebugLog("Add NN to VX" + fmt.Sprintf("(%X, %X)", VX, NN))
	case 0x8:
		switch N {
		case 0x0: // Set VX to VY
			c.Register[VX] = c.Register[VY]
			DebugLog("Set VX to VY")
		case 0x1: // Set VX to VX | VY
			c.Register[VX] |= c.Register[VY]
			DebugLog("Set VX to VX | VY")
		case 0x2: // Set VX to VX & VY
			c.Register[VX] &= c.Register[VY]
			DebugLog("Set VX to VX & VY")
		case 0x3: // Set VX to VX ^ VY
			c.Register[VX] ^= c.Register[VY]
			DebugLog("Set VX to VX ^ VY")
		case 0x4: // Add VY to VX and set VF to 1 if carry
			c.Register[0xF] = 0
			if (int)(c.Register[VX])+(int)(c.Register[VY]) > 0xFF {
				c.Register[0xF] = 1
			}
			c.Register[VX] += c.Register[VY]
			DebugLog("Add VY to VX and set VF to 1 if carry")
		case 0x5: // Subtract VY from VX and set VF to 0 if borrow
			c.Register[0xF] = 0
			if c.Register[VX] > c.Register[VY] {
				c.Register[0xF] = 1
			}
			c.Register[VX] -= c.Register[VY]
			DebugLog("Subtract VY from VX and set VF to 0 if borrow")
		case 0x6: // Shift VX right by one and set VF to least significant bit
			c.Register[0xF] = c.Register[VX] & 0x1
			c.Register[VX] >>= 1
			DebugLog("Shift VX right by one and set VF to least significant bit")
		case 0x7: // Subtract VY from VX and set VF to 0 if borrow
			c.Register[0xF] = 0
			if c.Register[VY] > c.Register[VX] {
				c.Register[0xF] = 1
			}
			c.Register[VX] = c.Register[VY] - c.Register[VX]
			DebugLog("Subtract VY from VX and set VF to 0 if borrow")
		case 0xE: // Shift VX left by one and set VF to most significant bit
			c.Register[0xF] = c.Register[VX] >> 7
			c.Register[VX] <<= 1
			DebugLog("Shift VX left by one and set VF to most significant bit")
		default:
			fmt.Printf("WARNING: Unknown opcode: 0x%X\n", N)
		}
	case 0x9: // Skip if VX != VY
		if c.Register[VX] != c.Register[VY] {
			c.Pc += 2
		}
		DebugLog("Skip if VX != VY")
	case 0xA: // Set IReg to NNN
		c.IReg = NNN
		DebugLog("Set IReg to NNN " + fmt.Sprintf("(%X, %X)", VX, NN))
	case 0xB: // Jump to NNN + V0
		c.Pc = NNN + uint16(c.Register[0])
		DebugLog("Jump to NNN + V0")
	case 0xC: // Set VX to random byte & NN
		c.Register[VX] = byte(rand.Intn(0x100)) & NN
		DebugLog("Set VX to random byte & NN")
	case 0xD: // Draw sprite at VX, VY with N bytes
		x := int(c.Register[VX] % 64)
		y := int(c.Register[VY] % 32)
		c.Register[0xF] = 0

		for i := 0; i < int(N); i++ { // row
			pixelY := y + i // Vertical position
			rowByte := c.Ram[int(c.IReg)+i]

			for j := 0; j < 8; j++ { // column/pixel
				pixelX := x + (7 - j) // Horizontal position
				pixelValue := (rowByte>>j)&0x1 == 1

				if c.Screen.FlipPixel(pixelX, pixelY, pixelValue) {
					c.Register[0xF] = 1
				}
			}
		}
		DebugLog("Draw sprite at VX, VY with N bytes")
	case 0xE: // Key input
		if c.Register[VX] > 0xF {
			fmt.Printf("WARNING: Unknown opcode: 0x%X\n", N)
			break
		}

		switch NN {
		case 0x9E: // Skip if key pressed
			if c.Pad.Keydata[c.Register[VX]] {
				c.Pc += 2
			}
			DebugLog("Skip if key pressed")
		case 0xA1: // Skip if key not pressed
			if !c.Pad.Keydata[c.Register[VX]] {
				c.Pc += 2
			}
			DebugLog("Skip if key not pressed")
		default:
			fmt.Printf("WARNING: Unknown opcode: 0x%X\n", NN)
		}
	case 0xF:
		switch NN {
		case 0x7: // Set VX to delay timer
			c.Register[VX] = c.Delay
			DebugLog("Set VX to delay timer")
		case 0xA: // Wait for key press
			c.Pc -= 2
			for i := byte(0); i < 0x10; i++ {
				if c.Pad.Keydata[i] {
					c.Pc += 2
					c.Register[VX] = 1
					break
				}
			}
			DebugLog("Wait for key press")
		case 0x15: // Set delay timer to VX
			c.Delay = c.Register[VX]
			DebugLog("Set delay timer to VX")
		case 0x18: // Set sound timer to VX
			// The Games has no sound
		case 0x1E: // Add VX to IReg
			c.IReg += uint16(c.Register[VX])
			DebugLog("Add VX to IReg")
		case 0x29: // Set IReg to sprite for digit VX
			//c.IReg = uint16(c.Register[VX]&0xF) * 5 // DEBUG
			c.IReg = uint16(VX) * 5 // DEBUG
			fmt.Print("WARNING: Not implemented opcode: 0xFx29\n")
		case 0x33: // Store BCD representation of VX in IReg
			c.Ram[c.IReg] = c.Register[VX] / 100
			c.Ram[c.IReg+1] = (c.Register[VX] / 10) % 10
			c.Ram[c.IReg+2] = (c.Register[VX] % 100) % 10 // DEBUG
			DebugLog("Store BCD representation of VX in IReg")
		case 0x55: // Store V0 through VX in IReg
			for i := uint16(0); i <= uint16(VX); i++ {
				c.Ram[c.IReg+i] = c.Register[i]
			}
			DebugLog("Store V0 through VX in IReg")
		case 0x65: // Load V0 through VX from IReg
			for i := uint16(0); i <= uint16(VX); i++ {
				c.Register[i] = c.Ram[c.IReg+i]
			}
			DebugLog("Load V0 through VX from IReg")
		default:
			fmt.Printf("WARNING: Unknown opcode: 0x%X\n", NN)
		}
	default:
		fmt.Printf("WARNING: Unknown opcode: 0x%X\n", N)
	}
}
