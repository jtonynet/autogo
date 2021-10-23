package handlers

import (
	"fmt"
	"log"

	"gobot.io/x/gobot/platforms/keyboard"

	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

//TODO env vars on viper
const (
	VERSION              = "v0.0.5"
	PAN_TILT_FACTOR      = 30
	MAX_SPEED            = 255
	MIN_STOP_SONAR_VALUE = 15
)

func InitKeyboard(keys *keyboard.Driver, motors *output.Motors, servoKit *output.Servos, SonarSet *input.SonarSet, lcd *output.Display) {
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
		key := data.(keyboard.KeyEvent)

		if key.Key == keyboard.B {
			sonarData, err := SonarSet.GetData()
			if err == nil {
				log.Println("///*********")
				log.Println("///Print arduino sonar data::")
				log.Println(sonarData)
				log.Println("///*********")
			}

		}

		panAngle := int(servoPan.CurrentAngle)
		tiltAngle := int(servoTilt.CurrentAngle)
		if key.Key == keyboard.W {
			newTilt := tiltAngle - PAN_TILT_FACTOR
			if newTilt < servoKit.TiltPos["top"] {
				newTilt = servoKit.TiltPos["top"]
			}
			servoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.S {
			newTilt := tiltAngle + PAN_TILT_FACTOR
			if newTilt > servoKit.TiltPos["down"] {
				newTilt = servoKit.TiltPos["down"]
			}
			servoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.A {
			newPan := panAngle + PAN_TILT_FACTOR
			if newPan > servoKit.PanPos["left"] {
				newPan = servoKit.PanPos["left"]
			}
			servoKit.SetAngle(servoPan, uint8(newPan))

		} else if key.Key == keyboard.D {
			newPan := panAngle - PAN_TILT_FACTOR
			if newPan < servoKit.PanPos["right"] {
				newPan = servoKit.PanPos["right"]
			}
			servoKit.SetAngle(servoPan, uint8(newPan))

		} else if key.Key == keyboard.X {
			servoKit.SetCenter(servoPan)
			servoKit.SetAngle(servoTilt, uint8(servoKit.TiltPos["horizon"]))
		}

		if key.Key == keyboard.ArrowUp {
			motors.Forward(MAX_SPEED)
			lcd.ShowMessage("Front", output.LINE_2)
		} else if key.Key == keyboard.ArrowDown {
			motors.Backward(MAX_SPEED)
			lcd.ShowMessage("Back", output.LINE_2)
		} else if key.Key == keyboard.ArrowRight {
			motors.Left(MAX_SPEED)
			lcd.ShowMessage("Left", output.LINE_2)
		} else if key.Key == keyboard.ArrowLeft {
			motors.Right(MAX_SPEED)
			lcd.ShowMessage("Right", output.LINE_2)
		} else if key.Key == keyboard.Q {
			motors.Stop()
			lcd.ShowMessage(VERSION+" Arrow key", output.LINE_2)
		} else {
			fmt.Println("keyboard event!", key, key.Char)
		}
	})
}
