package charLCDRGBI2C

import (
	"time"

	"github.com/googolgl/go-i2c"
	"github.com/googolgl/go-mcp23017"
)

const (
	// Registers
	IODIRA = 0x00 // I/O direction register for Port A
	IODIRB = 0x01

	// MCP23017 pin mappings based on Python library
	LcdRsPin     = "B7" // Pin 15
	LcdEnablePin = "B5" // Pin 13
	LcdD4Pin     = "B4" // Pin 12
	LcdD5Pin     = "B3" // Pin 11
	LcdD6Pin     = "B2" // Pin 10
	LcdD7Pin     = "B1" // Pin 9
	RwPin        = "B6" // Pin 14

	// MCP23017 pins for RGB LED
	RedPin       = "A6" // Pin 6
	GreenPin     = "A7" // Pin 7
	BluePin      = "B0" // Pin 8
	BacklightPin = "A5" // Pin 5

	// MCP23017 pins for Buttons
	LeftButton   = "A4" // Pin 4
	UpButton     = "A3" // Pin 3
	DownButton   = "A2" // Pin 2
	RightButton  = "A1" // Pin 1
	SelectButton = "A0" // Pin 0

	// Constants for LCD commands
	LCD_CLEARDISPLAY   = 0x01
	LCD_RETURNHOME     = 0x02
	LCD_ENTRYMODESET   = 0x04
	LCD_DISPLAYCONTROL = 0x08
	LCD_CURSORSHIFT    = 0x10
	LCD_FUNCTIONSET    = 0x20
	LCD_SETCGRAMADDR   = 0x40
	LCD_SETDDRAMADDR   = 0x80

	// Entry flags
	LCD_ENTRYLEFT           = 0x02
	LCD_ENTRYSHIFTDECREMENT = 0x00

	// Control flags
	LCD_DISPLAYON = 0x04
	LCD_CURSORON  = 0x02
	LCD_CURSOROFF = 0x00
	LCD_BLINKON   = 0x01
	LCD_BLINKOFF  = 0x00

	// Move flags
	LCD_DISPLAYMOVE = 0x08
	LCD_MOVERIGHT   = 0x04
	LCD_MOVELEFT    = 0x00

	// Function set flags
	LCD_4BITMODE = 0x00
	LCD_2LINE    = 0x08
	LCD_1LINE    = 0x00
	LCD_5X8DOTS  = 0x00

	// Direction constants
	LEFT_TO_RIGHT = 0
	RIGHT_TO_LEFT = 1
)

// Row offset addresses for different LCD lines
var LCD_ROW_OFFSETS = []byte{0x00, 0x40, 0x14, 0x54}

// CharLCDRGBI2C represents a character LCD with an RGB LED controlled via I2C.
type CharLCDRGBI2C struct {
	mcp        *mcp23017.MCP23017 // I2C expander
	columns    int                // Number of columns on the LCD
	lines      int                // Number of lines on the LCD
	backlight  bool               // Backlight status
	rgb        [3]string          // RGB pins
	colorValue [3]int             // RGB color values (0-100)

	// Display control
	displayControl  byte   // Control byte for display settings
	displayMode     byte   // Display mode settings
	displayFunction byte   // Display function settings
	row             int    // Current row position
	column          int    // Current column position
	columnAlign     bool   // Column alignment setting
	message         string // Message to be displayed
	direction       int    // LEFT_TO_RIGHT or RIGHT_TO_LEFT
}

func New(i2c *i2c.Options, columns, lines int) (*CharLCDRGBI2C, error) {
	// Initialize MCP23017
	mcp, err := mcp23017.New(i2c)
	if err != nil {
		return nil, err
	}

	lcd := &CharLCDRGBI2C{
		mcp:        mcp,
		columns:    columns,
		lines:      lines,
		backlight:  true,
		rgb:        [3]string{RedPin, GreenPin, BluePin},
		colorValue: [3]int{0, 0, 0},
	}

	lcd.setupPins()

	lcd.initialize()

	return lcd, nil
}

func (lcd *CharLCDRGBI2C) setupPins() {
	// Set LCD control pins as outputs
	lcd.mcp.Set(mcp23017.Pins{LcdRsPin, LcdEnablePin, LcdD4Pin, LcdD5Pin, LcdD6Pin, LcdD7Pin}).OUTPUT()
	lcd.mcp.Set(mcp23017.Pins{RwPin}).OUTPUT()

	// Set RGB LED pins as outputs
	lcd.mcp.Set(mcp23017.Pins{RedPin, GreenPin, BluePin}).OUTPUT()

	// Set Button pins as inputs with pull-up
	lcd.mcp.Set(mcp23017.Pins{LeftButton, UpButton, DownButton, RightButton, SelectButton}).INPUT()
	lcd.mcp.Set(mcp23017.Pins{LeftButton, UpButton, DownButton, RightButton, SelectButton}).PULLUP()
}

func (lcd *CharLCDRGBI2C) initialize() {
	// Wait for LCD to be ready
	time.Sleep(50 * time.Millisecond)

	// Pull RS low to begin commands
	lcd.mcp.Set(mcp23017.Pins{LcdRsPin}).LOW()
	lcd.mcp.Set(mcp23017.Pins{LcdEnablePin}).LOW()
	lcd.mcp.Set(mcp23017.Pins{RwPin}).LOW() // Write mode

	// 4-bit mode initialization sequence
	lcd.write4bits(0x03)
	time.Sleep(5 * time.Millisecond)
	lcd.write4bits(0x03)
	time.Sleep(5 * time.Millisecond)
	lcd.write4bits(0x03)
	time.Sleep(1 * time.Millisecond)
	lcd.write4bits(0x02) // Set to 4-bit mode
	time.Sleep(1 * time.Millisecond)

	// Initialize display control
	lcd.displayControl = LCD_DISPLAYON | LCD_CURSOROFF | LCD_BLINKOFF
	lcd.displayFunction = LCD_4BITMODE | LCD_1LINE | LCD_2LINE | LCD_5X8DOTS
	lcd.displayMode = LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT

	// Write to displaycontrol
	lcd.write8(LCD_DISPLAYCONTROL | lcd.displayControl)
	// Write to displayfunction
	lcd.write8(LCD_FUNCTIONSET | lcd.displayFunction)
	// Set entry mode
	lcd.write8(LCD_ENTRYMODESET | lcd.displayMode)

	// Clear display
	lcd.Clear()

	// Initialize tracking variables
	lcd.row = 0
	lcd.column = 0
	lcd.columnAlign = false
	lcd.direction = LEFT_TO_RIGHT
	lcd.message = ""

	// Turn off all RGB LEDs initially
	lcd.SetColor(0, 0, 0)
}

// Clear clears the LCD display
func (lcd *CharLCDRGBI2C) Clear() {
	lcd.write8(LCD_CLEARDISPLAY)
	time.Sleep(3 * time.Millisecond) // This command takes a long time
}

// Home moves cursor to home position
func (lcd *CharLCDRGBI2C) Home() {
	lcd.write8(LCD_RETURNHOME)
	time.Sleep(3 * time.Millisecond) // This command takes a long time
}

// CursorPosition sets the cursor position
func (lcd *CharLCDRGBI2C) CursorPosition(column, row int) {
	// Clamp row to the last row of the display
	if row >= lcd.lines {
		row = lcd.lines - 1
	}
	// Clamp to last column of display
	if column >= lcd.columns {
		column = lcd.columns - 1
	}
	// Set location
	lcd.write8(LCD_SETDDRAMADDR | (byte(column) + LCD_ROW_OFFSETS[row]))
	// Update row and column tracking
	lcd.row = row
	lcd.column = column
}

// SetCursor enables or disables the cursor
func (lcd *CharLCDRGBI2C) SetCursor(show bool) {
	if show {
		lcd.displayControl |= LCD_CURSORON
	} else {
		lcd.displayControl &= ^byte(LCD_CURSORON) // Use explicit type conversion
	}
	lcd.write8(LCD_DISPLAYCONTROL | lcd.displayControl)
}

// SetBlink enables or disables cursor blinking
func (lcd *CharLCDRGBI2C) SetBlink(blink bool) {
	if blink {
		lcd.displayControl |= LCD_BLINKON
	} else {
		lcd.displayControl &= ^byte(LCD_BLINKON) // Use explicit type conversion
	}
	lcd.write8(LCD_DISPLAYCONTROL | lcd.displayControl)
}

// SetDisplay enables or disables the entire display
func (lcd *CharLCDRGBI2C) SetDisplay(enable bool) {
	if enable {
		lcd.displayControl |= LCD_DISPLAYON
	} else {
		lcd.displayControl &= ^byte(LCD_DISPLAYON) // Use explicit type conversion
	}
	lcd.write8(LCD_DISPLAYCONTROL | lcd.displayControl)
}

// MoveLeft moves displayed text left one column
func (lcd *CharLCDRGBI2C) MoveLeft() {
	lcd.write8(LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVELEFT)
}

// MoveRight moves displayed text right one column
func (lcd *CharLCDRGBI2C) MoveRight() {
	lcd.write8(LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVERIGHT)
}

// SetTextDirection sets the text direction
func (lcd *CharLCDRGBI2C) SetTextDirection(direction int) {
	lcd.direction = direction
	if direction == LEFT_TO_RIGHT {
		lcd.leftToRight()
	} else {
		lcd.rightToLeft()
	}
}

// leftToRight sets text direction from left to right
func (lcd *CharLCDRGBI2C) leftToRight() {
	lcd.displayMode |= LCD_ENTRYLEFT
	lcd.write8(LCD_ENTRYMODESET | lcd.displayMode)
}

// rightToLeft sets text direction from right to left
func (lcd *CharLCDRGBI2C) rightToLeft() {
	lcd.displayMode &= ^byte(LCD_ENTRYLEFT) // Use explicit type conversion
	lcd.write8(LCD_ENTRYMODESET | lcd.displayMode)
}

// SetColumnAlign sets column alignment for newlines
func (lcd *CharLCDRGBI2C) SetColumnAlign(enable bool) {
	lcd.columnAlign = enable
}

// CreateChar creates a custom character
func (lcd *CharLCDRGBI2C) CreateChar(location byte, pattern []byte) {
	// Only positions 0-7 are allowed
	location &= 0x7
	lcd.write8(LCD_SETCGRAMADDR | (location << 3))
	for i := 0; i < 8; i++ {
		lcd.write8(pattern[i], true)
	}
}

// Message displays text on the LCD
func (lcd *CharLCDRGBI2C) Message(message string) {
	lcd.message = message

	// Set line to match current row
	line := lcd.row
	// Track initial character
	initialCharacter := 0

	// Iterate through each character
	for _, character := range message {
		// If this is the first character in the string
		if initialCharacter == 0 {
			// Start at current position determined by text direction
			var col int
			if lcd.displayMode&LCD_ENTRYLEFT > 0 {
				col = lcd.column
			} else {
				col = lcd.columns - 1 - lcd.column
			}
			lcd.CursorPosition(col, line)
			initialCharacter++
		}

		// If character is newline, go to next line
		if character == '\n' {
			line++
			// Handle starting position on new line
			var col int
			if lcd.displayMode&LCD_ENTRYLEFT > 0 {
				if lcd.columnAlign {
					col = lcd.column
				} else {
					col = 0
				}
			} else {
				if lcd.columnAlign {
					col = lcd.column
				} else {
					col = lcd.columns - 1
				}
			}
			lcd.CursorPosition(col, line)
		} else {
			// Write character to display
			lcd.write8(byte(character), true)
		}
	}

	// Reset column and row to (0,0) after message is displayed
	lcd.column, lcd.row = 0, 0
}

// write8 sends 8-bit value to the LCD
func (lcd *CharLCDRGBI2C) write8(value byte, charMode ...bool) {
	// Default to command mode (false)
	isCharMode := false
	if len(charMode) > 0 {
		isCharMode = charMode[0]
	}

	// Set RS pin based on character/command mode
	if isCharMode {
		lcd.mcp.Set(mcp23017.Pins{LcdRsPin}).HIGH() // Character mode
	} else {
		lcd.mcp.Set(mcp23017.Pins{LcdRsPin}).LOW() // Command mode
	}

	// Write upper 4 bits
	lcd.write4bits(value >> 4)
	// Write lower 4 bits
	lcd.write4bits(value & 0x0F)
}

// write4bits sends 4-bits to the LCD
func (lcd *CharLCDRGBI2C) write4bits(value byte) {
	// Set data pins
	if value&0x01 > 0 {
		lcd.mcp.Set(mcp23017.Pins{LcdD4Pin}).HIGH()
	} else {
		lcd.mcp.Set(mcp23017.Pins{LcdD4Pin}).LOW()
	}

	if value&0x02 > 0 {
		lcd.mcp.Set(mcp23017.Pins{LcdD5Pin}).HIGH()
	} else {
		lcd.mcp.Set(mcp23017.Pins{LcdD5Pin}).LOW()
	}

	if value&0x04 > 0 {
		lcd.mcp.Set(mcp23017.Pins{LcdD6Pin}).HIGH()
	} else {
		lcd.mcp.Set(mcp23017.Pins{LcdD6Pin}).LOW()
	}

	if value&0x08 > 0 {
		lcd.mcp.Set(mcp23017.Pins{LcdD7Pin}).HIGH()
	} else {
		lcd.mcp.Set(mcp23017.Pins{LcdD7Pin}).LOW()
	}

	// Pulse enable pin
	lcd.pulseEnable()
}

// pulseEnable pulses the enable pin to latch command
func (lcd *CharLCDRGBI2C) pulseEnable() {
	lcd.mcp.Set(mcp23017.Pins{LcdEnablePin}).LOW()
	time.Sleep(1 * time.Microsecond)
	lcd.mcp.Set(mcp23017.Pins{LcdEnablePin}).HIGH()
	time.Sleep(1 * time.Microsecond)
	lcd.mcp.Set(mcp23017.Pins{LcdEnablePin}).LOW()
	time.Sleep(100 * time.Microsecond) // Commands need > 37us to settle
}
