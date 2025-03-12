package charLCDRGBI2C

import (
	"log"

	"github.com/googolgl/go-mcp23017"
)

// Backlight
func (lcd *CharLCDRGBI2C) SetBacklight(on bool) error {
	if on {
		// Set as output to turn backlight ON
		lcd.mcp.Set(mcp23017.Pins{BacklightPin}).OUTPUT()
		log.Println("Backlight ON")
	} else {
		// Set as input to turn backlight OFF
		lcd.mcp.Set(mcp23017.Pins{BacklightPin}).INPUT()
		log.Println("Backlight OFF")
	}
	return nil
}
