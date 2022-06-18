package peripherals

import (
	"gobot.io/x/gobot/platforms/keyboard"
)

/*
KeyEvent overwriten from keyboard here to maintence
implementation keybord driver outside of application
*/
type KeyEvent struct {
	Key  int
	Char string
}

type Keyboard struct {
	Driver *keyboard.Driver
	Key    string
}

func GetKeyboard() *Keyboard {
	this := &Keyboard{Driver: keyboard.NewDriver(), Key: keyboard.Key}
	return this
}

func GetKeyEvent(data interface{}) KeyEvent {
	k := data.(keyboard.KeyEvent)

	var key KeyEvent
	key.Char = k.Char
	key.Key = k.Key

	return key
}
