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

	Button(lcd)
}

func Button(lcd *charLCDRGBI2C.CharLCDRGBI2C) {
	log.Println("Starting Button Demo")

	debounceTime := 200 * time.Millisecond
	lastPressTime := time.Now()

	for {
		// Simple debouncing
		if time.Since(lastPressTime) < debounceTime {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		var buttonMessage string

		switch {
		case lcd.LeftButton():
			buttonMessage = "Left"
			lastPressTime = time.Now()
		case lcd.UpButton():
			buttonMessage = "Up"
			lastPressTime = time.Now()
		case lcd.DownButton():
			buttonMessage = "Down"
			lastPressTime = time.Now()
		case lcd.RightButton():
			buttonMessage = "Right"
			lastPressTime = time.Now()
		case lcd.SelectButton():
			buttonMessage = "Select"
			lastPressTime = time.Now()
		default:
			// No button pressed, continue checking
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// Wait a bit to show the result
		time.Sleep(500 * time.Millisecond)

		log.Printf("Button pressed: %s", buttonMessage)

		// Give time for button release
		time.Sleep(debounceTime)
	}
}
