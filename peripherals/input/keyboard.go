package peripherals

import (
	"gobot.io/x/gobot/platforms/keyboard"
)

type Keyboard struct {
	Driver *keyboard.Driver
	Key    string
	//KeyEvent    *keyboard.KeyEvent
}

func GetKeyboard() *Keyboard {
	//TODO: dont pass type keyboard.KeyEvent into type input.Keyboard, wrappers fail :(
	//fix keyboard adapter in future to remove gobot/platforms/keyboard from application and domain
	this := &Keyboard{Driver: keyboard.NewDriver(), Key: keyboard.Key}
	return this
}
