// FROM: https://github.com/adafruit/Adafruit_CircuitPython_CharLCD/blob/main/examples/charlcd_customcharacter.py
package main

import (
	"log"

	"github.com/googolgl/go-i2c"
	"github.com/googolgl/go-mcp23017"
	"github.com/jyap808/charLCDRGBI2C"
)

func main() {
	// Initialize I2C
	i2c, err := i2c.New(mcp23017.DefI2CAdr, "/dev/i2c-1")
	if err != nil {
		log.Fatalf("Failed to initialize I2C: %v", err)
	}
	defer i2c.Close()

	// Create LCD object (16 columns, 2 rows)
	lcd, err := charLCDRGBI2C.New(i2c, 16, 2)
	if err != nil {
		log.Fatalf("Failed to initialize LCD: %v", err)
	}

	CustomCharacter(lcd)
}

func CustomCharacter(lcd *charLCDRGBI2C.CharLCDRGBI2C) {
	checkmark := []byte{0x0, 0x0, 0x1, 0x3, 0x16, 0x1C, 0x8, 0x0}

	// Store in LCD character memory 0
	lcd.CreateChar(0, checkmark)

	lcd.Clear()
	lcd.Message("\x00 Success \x00")
}
