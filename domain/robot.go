package domain

import (
	"fmt"
	"log"
	"net"

	config "github.com/jtonynet/autogo/config"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
	"gobot.io/x/gobot/platforms/keyboard"
)

var (
	direction         string = ""
	LCDMsg            string = ""
	colissionDetected bool   = false
)

type Robot struct {
	Motors   *output.Motors
	ServoKit *output.Servos
	LCD      *output.Display
	SonarSet *input.SonarSet
	Cfg      *config.Config
}

func NewRobot(Motors *output.Motors, ServoKit *output.Servos, LCD *output.Display, SonarSet *input.SonarSet, Cfg *config.Config) *Robot {
	this := &Robot{Motors: Motors, ServoKit: ServoKit, LCD: LCD, SonarSet: SonarSet, Cfg: Cfg}

	if Cfg.ServoKit.Enabled {
		servoPan := ServoKit.GetByName("pan")
		servoTilt := ServoKit.GetByName("tilt")

		ServoKit.Init()
		ServoKit.SetCenter(servoPan)
		ServoKit.SetAngle(servoTilt, uint8(ServoKit.TiltPos["horizon"]))
	}

	if Cfg.LCD.Enabled {
		ip := getOutboundIP()
		err := LCD.ShowMessage(string(ip), output.LINE_1)
		if err != nil {
			log.Fatal(err)
		}

		err = LCD.ShowMessage(Cfg.Version+" Arrow key", output.LINE_2)
		if err != nil {
			log.Fatal(err)
		}
	}

	if Cfg.ArduinoSonar.Enabled && Cfg.Motors.Enabled {
		go this.sonarWorker()
	}

	return this
}

func (this *Robot) ControllByKeyboard(data interface{}) {
	oldDirection := direction
	key := input.GetKeyEvent(data)

	if this.Cfg.ServoKit.Enabled {
		servoPan := this.ServoKit.GetByName("pan")
		servoTilt := this.ServoKit.GetByName("tilt")

		panAngle := int(servoPan.CurrentAngle)
		tiltAngle := int(servoTilt.CurrentAngle)

		if key.Key == keyboard.W {
			newTilt := tiltAngle - this.Cfg.ServoKit.PanTiltFactor
			if newTilt < this.ServoKit.TiltPos["top"] {
				newTilt = this.ServoKit.TiltPos["top"]
			}
			this.ServoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.S {
			newTilt := tiltAngle + this.Cfg.ServoKit.PanTiltFactor
			if newTilt > this.ServoKit.TiltPos["down"] {
				newTilt = this.ServoKit.TiltPos["down"]
			}
			this.ServoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.A {
			newPan := panAngle + this.Cfg.ServoKit.PanTiltFactor
			if newPan > this.ServoKit.PanPos["left"] {
				newPan = this.ServoKit.PanPos["left"]
			}
			this.ServoKit.SetAngle(servoPan, uint8(newPan))

		} else if key.Key == keyboard.D {
			newPan := panAngle - this.Cfg.ServoKit.PanTiltFactor
			if newPan < this.ServoKit.PanPos["right"] {
				newPan = this.ServoKit.PanPos["right"]
			}
			this.ServoKit.SetAngle(servoPan, uint8(newPan))
		} else if key.Key == keyboard.X {
			this.ServoKit.SetCenter(servoPan)
			this.ServoKit.SetAngle(servoTilt, uint8(this.ServoKit.TiltPos["horizon"]))
		}
	}

	if this.Cfg.Motors.Enabled {
		if key.Key == keyboard.ArrowUp && colissionDetected == false {
			this.Motors.Forward(this.Cfg.Motors.MaxSpeed)
			direction = "Front"
			LCDMsg = direction
		} else if key.Key == keyboard.ArrowDown {
			this.Motors.Backward(this.Cfg.Motors.MaxSpeed)
			direction = "Back"
			LCDMsg = direction
		} else if key.Key == keyboard.ArrowRight {
			this.Motors.Left(this.Cfg.Motors.MaxSpeed)
			direction = "Right"
			LCDMsg = direction
		} else if key.Key == keyboard.ArrowLeft {
			this.Motors.Right(this.Cfg.Motors.MaxSpeed)
			direction = "Left"
			LCDMsg = direction
		} else if key.Key == keyboard.Q {
			this.Motors.Stop()
			direction = ""
			LCDMsg = this.Cfg.Version + " Arrow key"
		} else {
			fmt.Println(LCDMsg, key, key.Char)
		}
	}

	if this.Cfg.LCD.Enabled && oldDirection != direction {
		this.LCD.ShowMessage(LCDMsg, output.LINE_2)
	}
}

func (this *Robot) sonarWorker() {
	for true {
		sonarData, err := this.SonarSet.GetData()
		if err == nil {
			if sonarData["center"] <= this.Cfg.ArduinoSonar.MinStopValue && direction == "Front" && colissionDetected == false {
				colissionDetected = true
				this.Motors.Stop()

				if this.Cfg.LCD.Enabled {
					s := fmt.Sprintf("STOP CRASH %.2f", sonarData["center"])
					this.LCD.ShowMessage(s, output.LINE_2)
				}

			} else if colissionDetected && direction != "Front" {
				colissionDetected = false
			}
		}
	}
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "ip offline"
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
