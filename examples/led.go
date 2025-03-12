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

	LED(lcd)
}

func LED(lcd *charLCDRGBI2C.CharLCDRGBI2C) {
	log.Println("Starting RGB LED Demo")

	// Cycle through some colors
	colorMap := map[string][3]int{
		"Red":         {100, 0, 0},
		"Green":       {0, 100, 0},
		"Blue":        {0, 0, 100},
		"Purple":      {50, 0, 50},
		"Cyan":        {0, 50, 50},
		"Yellow":      {50, 50, 0},
		"White (dim)": {50, 50, 50},
	}

	for color, colorV := range colorMap {
		log.Printf("Setting color to: R=%d, G=%d, B=%d - %s", colorV[0], colorV[1], colorV[2], color)
		lcd.SetColor(colorV[0], colorV[1], colorV[2])
		time.Sleep(1 * time.Second)
	}

	// Turn off all LEDs
	lcd.SetColor(0, 0, 0)
}
