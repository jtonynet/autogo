package application

import (
	"fmt"

	"gobot.io/x/gobot/platforms/keyboard"

	"github.com/jtonynet/autogo/config"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

var (
	direction         string = ""
	colissionDetected bool   = false
)

func InitKeyboard(keys *keyboard.Driver, motors *output.Motors, servoKit *output.Servos, sonarSet *input.SonarSet, lcd *output.Display, config config.Config) {
	firstRun := 1
	servoPan := servoKit.GetByName("pan")
	servoTilt := servoKit.GetByName("tilt")

	if firstRun == 1 {
		firstRun = 0
		servoKit.Init()
		servoKit.SetCenter(servoPan)
		servoKit.SetAngle(servoTilt, uint8(servoKit.TiltPos["horizon"]))
	}

	keys.On(keyboard.Key, func(data interface{}) {
		go sonarWorker(sonarSet, motors, lcd, config)

		key := data.(keyboard.KeyEvent)

		panAngle := int(servoPan.CurrentAngle)
		tiltAngle := int(servoTilt.CurrentAngle)
		if key.Key == keyboard.W {
			newTilt := tiltAngle - config.PanTiltFactor
			if newTilt < servoKit.TiltPos["top"] {
				newTilt = servoKit.TiltPos["top"]
			}
			servoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.S {
			newTilt := tiltAngle + config.PanTiltFactor
			if newTilt > servoKit.TiltPos["down"] {
				newTilt = servoKit.TiltPos["down"]
			}
			servoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.A {
			newPan := panAngle + config.PanTiltFactor
			if newPan > servoKit.PanPos["left"] {
				newPan = servoKit.PanPos["left"]
			}
			servoKit.SetAngle(servoPan, uint8(newPan))

		} else if key.Key == keyboard.D {
			newPan := panAngle - config.PanTiltFactor
			if newPan < servoKit.PanPos["right"] {
				newPan = servoKit.PanPos["right"]
			}
			servoKit.SetAngle(servoPan, uint8(newPan))

		} else if key.Key == keyboard.X {
			servoKit.SetCenter(servoPan)
			servoKit.SetAngle(servoTilt, uint8(servoKit.TiltPos["horizon"]))
		}

		if key.Key == keyboard.ArrowUp && colissionDetected == false {
			motors.Forward(config.MaxMotorsSpeed)
			direction = "Front"
			lcd.ShowMessage(direction, output.LINE_2)
		} else if key.Key == keyboard.ArrowDown {
			motors.Backward(config.MaxMotorsSpeed)
			direction = "Back"
			lcd.ShowMessage(direction, output.LINE_2)
		} else if key.Key == keyboard.ArrowRight {
			motors.Left(config.MaxMotorsSpeed)
			direction = "Right"
			lcd.ShowMessage(direction, output.LINE_2)
		} else if key.Key == keyboard.ArrowLeft {
			motors.Right(config.MaxMotorsSpeed)
			direction = "Left"
			lcd.ShowMessage(direction, output.LINE_2)
		} else if key.Key == keyboard.Q {
			motors.Stop()
			direction = ""
			lcd.ShowMessage(config.Version+" Arrow key", output.LINE_2)
		} else {
			fmt.Println("keyboard event!", key, key.Char)
		}
	})
}

func sonarWorker(sonarSet *input.SonarSet, motors *output.Motors, lcd *output.Display, config config.Config) {
	for true {
		sonarData, err := sonarSet.GetData()
		if err == nil {
			if config.MinStopSonarValue >= sonarData["center"] && colissionDetected == false && direction == "Front" {
				colissionDetected = true
				motors.Stop()

				s := fmt.Sprintf("STOP CRASH %.2f", sonarData["center"])
				lcd.ShowMessage(s, output.LINE_2)
			} else if colissionDetected && direction != "Front" {
				colissionDetected = false
			}
		}
	}
}