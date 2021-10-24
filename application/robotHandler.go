package application

import (
	"fmt"

	"gobot.io/x/gobot/platforms/keyboard"

	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

//TODO env vars on viper
const (
	VERSION                      = "v0.0.3"
	PAN_TILT_FACTOR              = 30
	MAX_SPEED                    = 255
	MIN_STOP_SONAR_VALUE float64 = 15.00
)

var (
	direction         string = ""
	colissionDetected bool   = false
)

func InitKeyboard(keys *keyboard.Driver, motors *output.Motors, servoKit *output.Servos, sonarSet *input.SonarSet, lcd *output.Display) {
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
		go sonarWorker(sonarSet, motors, lcd)

		key := data.(keyboard.KeyEvent)

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

		if key.Key == keyboard.ArrowUp && colissionDetected == false {
			motors.Forward(MAX_SPEED)
			direction = "Front"
			lcd.ShowMessage(direction, output.LINE_2)
		} else if key.Key == keyboard.ArrowDown {
			motors.Backward(MAX_SPEED)
			direction = "Back"
			lcd.ShowMessage(direction, output.LINE_2)
		} else if key.Key == keyboard.ArrowRight {
			motors.Left(MAX_SPEED)
			direction = "Right"
			lcd.ShowMessage(direction, output.LINE_2)
		} else if key.Key == keyboard.ArrowLeft {
			motors.Right(MAX_SPEED)
			direction = "Left"
			lcd.ShowMessage(direction, output.LINE_2)
		} else if key.Key == keyboard.Q {
			motors.Stop()
			direction = ""
			lcd.ShowMessage(VERSION+" Arrow key", output.LINE_2)
		} else {
			fmt.Println("keyboard event!", key, key.Char)
		}
	})
}

func sonarWorker(sonarSet *input.SonarSet, motors *output.Motors, lcd *output.Display) {
	for true {
		sonarData, err := sonarSet.GetData()
		if err == nil {
			if MIN_STOP_SONAR_VALUE >= sonarData["center"] && colissionDetected == false && direction == "Front" {
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
