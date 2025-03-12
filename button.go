package charLCDRGBI2C

import (
	"log"

	"github.com/googolgl/go-mcp23017"
)

// IsButtonPressed checks if a specific button is pressed
func (lcd *CharLCDRGBI2C) IsButtonPressed(buttonPin string) bool {
	// Read the button state (LOW when pressed because of pull-up resistor)
	pinStates, err := lcd.mcp.Get(mcp23017.Pins{buttonPin})
	if err != nil {
		log.Printf("Error reading button state: %v", err)
		return false
	}

	// Check if the button's value is in the map and is LOW (pressed)
	value, exists := pinStates[buttonPin]
	if !exists {
		log.Printf("Button pin %s not found in state map", buttonPin)
		return false
	}

	// Return true if button is pressed (LOW)
	return value == 0 // 0 means LOW which means pressed (due to pull-up)
}

// Button state properties
func (lcd *CharLCDRGBI2C) LeftButton() bool {
	return lcd.IsButtonPressed(LeftButton)
}

func (lcd *CharLCDRGBI2C) UpButton() bool {
	return lcd.IsButtonPressed(UpButton)
}

func (lcd *CharLCDRGBI2C) DownButton() bool {
	return lcd.IsButtonPressed(DownButton)
}

func (lcd *CharLCDRGBI2C) RightButton() bool {
	return lcd.IsButtonPressed(RightButton)
}

func (lcd *CharLCDRGBI2C) SelectButton() bool {
	return lcd.IsButtonPressed(SelectButton)
}
