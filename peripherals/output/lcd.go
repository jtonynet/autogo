package peripherals

import (
	"log"
	"strings"
	"time"

	deviceD2r2 "github.com/d2r2/go-hd44780"
	i2cD2r2 "github.com/d2r2/go-i2c"
	config "github.com/jtonynet/autogo/config"
)

const DEFAULT_COLLUMNS = 16

type Display struct {
	i2c      *i2cD2r2.I2C
	lcd      *deviceD2r2.Lcd
	collumns int
	InUse    bool
}

func NewLcd(cfg config.LCD) (*Display, error) {
	//TODO: Use lcd i2c gobot solution to 16x2 screen

	bus := cfg.Bus
	addr := cfg.Addr
	collumns := cfg.Collumns

	// Create new connection to i2c-bus on 2 line with address 0x27.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2cD2r2.NewI2C(addr, bus)
	if err != nil {
		log.Fatal(err)
	}

	// Construct lcd-device connected via I2C connection
	lcd, err := deviceD2r2.NewLcd(i2c, deviceD2r2.LCD_16x2)
	if err != nil {
		log.Fatal(err)
	}

	if collumns == 0 {
		collumns = DEFAULT_COLLUMNS
	}

	// Turn on the backlight
	err = lcd.BacklightOn()
	if err != nil {
		log.Fatal(err)
	}

	this := &Display{i2c: i2c, lcd: lcd, collumns: collumns, InUse: false}
	return this, nil
}

func (this *Display) DeferAction() {
	defer this.i2c.Close()
}

func (this *Display) ShowMessage(message string, l int) error {
	time.Sleep(5 * time.Millisecond)

	lines := map[int]deviceD2r2.ShowOptions{
		1: deviceD2r2.SHOW_LINE_1,
		2: deviceD2r2.SHOW_LINE_2,
		3: deviceD2r2.SHOW_LINE_3,
		4: deviceD2r2.SHOW_LINE_4,
	}

	if this.InUse == false {
		this.InUse = true
		err := this.lcd.ShowMessage(rightPad(message, " ", this.collumns), lines[l])
		if err != nil {
			this.InUse = false
			return err
		}

		this.InUse = false
		return nil
	}

	time.Sleep(5 * time.Millisecond)
	return this.ShowMessage(message, l)
}

func rightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, (pLen-len(s)))
}
