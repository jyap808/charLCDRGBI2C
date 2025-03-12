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

	Message(lcd)
}

func Message(lcd *charLCDRGBI2C.CharLCDRGBI2C) {
	log.Println("Starting Message Demo")

	// Display a simple message
	lcd.Message("Hello, World!")
	time.Sleep(2 * time.Second)

	// Display a multi-line message
	lcd.Clear()
	lcd.Message("Line 1\nLine 2")
	time.Sleep(2 * time.Second)

	// Test cursor positioning
	lcd.Clear()
	lcd.CursorPosition(5, 0)
	lcd.Message("Position")
	time.Sleep(2 * time.Second)

	// Test scrolling
	lcd.Clear()
	lcd.Message("Scrolling text")
	time.Sleep(1 * time.Second)
	for i := 0; i < 5; i++ {
		lcd.MoveLeft()
		time.Sleep(500 * time.Millisecond)
	}
	lcd.Clear()
}
