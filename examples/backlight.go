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

	Backlight(lcd)
}

func Backlight(lcd *charLCDRGBI2C.CharLCDRGBI2C) {
	log.Println("Starting Backlight Demo")

	log.Println("Turning backlight OFF")
	lcd.SetBacklight(false)
	time.Sleep(2 * time.Second)

	lcd.Message("Backlight ON")
	log.Println("Turning backlight ON")
	lcd.SetBacklight(true)
	time.Sleep(2 * time.Second)

	log.Println("Turning backlight OFF again")
	lcd.SetBacklight(false)
	lcd.Clear()
}
