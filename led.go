package charLCDRGBI2C

import (
	"log"

	"github.com/googolgl/go-mcp23017"
)

// SetColor sets the RGB LED color (values from 0-100)
func (lcd *CharLCDRGBI2C) SetColor(red, green, blue int) {
	lcd.colorValue = [3]int{red, green, blue}

	// We need to invert the values as the Python code does (map 0-100 to on/off)
	// In Python, higher values = lower duty cycle, meaning 0=fully on, 100=fully off
	// We'll simulate this with digital pins

	// Update each LED
	values := [3]int{red, green, blue}
	pins := [3]string{RedPin, GreenPin, BluePin}

	for i, value := range values {
		if value > 1 {
			// Any value > 1 turns LED on (inverse of Python logic)
			lcd.mcp.Set(mcp23017.Pins{pins[i]}).LOW() // LOW = on for common anode RGB LED
		} else {
			lcd.mcp.Set(mcp23017.Pins{pins[i]}).HIGH() // HIGH = off
		}
	}
}

// SetColorRGB sets the RGB LED color using a 24-bit RGB integer
func (lcd *CharLCDRGBI2C) SetColorRGB(colorInt int) {
	if colorInt>>24 != 0 {
		log.Fatal("Integer color value must be positive and 24 bits max")
	}

	// Extract RGB components and convert to 0-100 scale
	r := float64(colorInt>>16) / 2.55
	g := float64((colorInt>>8)&0xFF) / 2.55
	b := float64(colorInt&0xFF) / 2.55

	lcd.SetColor(int(r), int(g), int(b))
}
