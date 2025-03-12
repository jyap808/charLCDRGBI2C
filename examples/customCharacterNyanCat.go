// FROM: https://github.com/adafruit/Adafruit_CircuitPython_CharLCD/blob/main/examples/charlcd_custom_character_nyan_cat.py
package main

import (
	"log"
	"time"

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

	CustomCharacterNyanCat(lcd)
}

func CustomCharacterNyanCat(lcd *charLCDRGBI2C.CharLCDRGBI2C) {
	head := []byte{31, 17, 27, 17, 17, 21, 17, 31}

	topBody := []byte{31, 0, 31, 0, 18, 8, 2, 8}
	topLeftCornerBody := []byte{31, 16, 16, 17, 22, 20, 20, 20}
	topRightCornerBody := []byte{31, 1, 1, 17, 13, 5, 5, 5}

	botBody := make([]byte, len(topBody))
	botLeftCornerBody := make([]byte, len(topLeftCornerBody))
	botRightCornerBody := make([]byte, len(topRightCornerBody))

	tailNeutral := []byte{0, 0, 0, 0, 31, 31, 0, 0}
	tailUp := []byte{0, 8, 12, 6, 3, 1, 0, 0}

	for i := 0; i < len(topBody); i++ {
		botBody[i] = topBody[len(topBody)-1-i]
		botLeftCornerBody[i] = topLeftCornerBody[len(topLeftCornerBody)-1-i]
		botRightCornerBody[i] = topRightCornerBody[len(topRightCornerBody)-1-i]
	}

	// Adding feet and making space for them
	botBody[6] = 31
	botBody[5] = 0
	botBody[4] = 31
	botBody[7] = 24
	botLeftCornerBody[7] = 0
	botLeftCornerBody[6] = 31
	botLeftCornerBody[7] = 28
	botRightCornerBody[7] = 0
	botRightCornerBody[6] = 31

	// Bottom body with feet forward
	botBody2 := append(botBody[:len(botBody)-1], 3)

	rainbow := []byte{0, 0, 6, 25, 11, 29, 27, 12}
	rainbow2 := []byte{0, 0, 6, 31, 13, 5, 23, 12}

	lcd.CreateChar(0, topBody)
	lcd.CreateChar(1, topLeftCornerBody)
	lcd.CreateChar(2, rainbow)
	lcd.CreateChar(3, botLeftCornerBody)
	lcd.CreateChar(4, botBody)
	lcd.CreateChar(5, botRightCornerBody)
	lcd.CreateChar(6, head)
	lcd.CreateChar(7, tailNeutral)

	lcd.Clear()
	lcd.MoveRight()
	lcd.Message("\x02\x02\x02\x02\x01\x00\x00\x00\x06\n\x02\x02\x02\x07\x03\x04\x04\x04\x05")

	lcd.SetBacklight(true)

	for {
		lcd.CreateChar(4, botBody2)
		lcd.CreateChar(7, tailUp)
		lcd.CreateChar(2, rainbow2)
		lcd.MoveRight()
		time.Sleep(400 * time.Millisecond)
		lcd.CreateChar(4, botBody)
		lcd.CreateChar(7, tailNeutral)
		lcd.CreateChar(2, rainbow)
		lcd.MoveLeft()
		time.Sleep(400 * time.Millisecond)
	}
}
