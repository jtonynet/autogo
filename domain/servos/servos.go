package domain

import (
	actuators "github.com/jtonynet/autogo/peripherals/actuators"
	"gobot.io/x/gobot/drivers/gpio"
)

type Servos struct {
	Kit  *actuators.Servos
	Pan  *gpio.ServoDriver
	Tilt *gpio.ServoDriver
}

func NewServos(servoKit *actuators.Servos) *Servos {
	this := &Servos{
		Kit:  servoKit,
		Pan:  servoKit.GetByName("pan"),
		Tilt: servoKit.GetByName("tilt"),
	}

	servoKit.Init()
	servoKit.SetCenter(this.Pan)
	servoKit.SetAngle(this.Tilt, uint8(servoKit.TiltPos["horizon"]))

	return this
}

func (this *Servos) ControllPanAndTilt(camDirection string) {
	cfg := this.Kit.Cfg
	servoPan := this.Pan
	servoTilt := this.Tilt

	panAngle := int(servoPan.CurrentAngle)
	tiltAngle := int(servoTilt.CurrentAngle)

	switch camDirection {
	case "Top":
		newTilt := tiltAngle - cfg.PanTiltFactor
		if newTilt < this.Kit.TiltPos["top"] {
			newTilt = this.Kit.TiltPos["top"]
		}
		this.Kit.SetAngle(servoTilt, uint8(newTilt))

	case "Down":
		newTilt := tiltAngle + cfg.PanTiltFactor
		if newTilt > this.Kit.TiltPos["down"] {
			newTilt = this.Kit.TiltPos["down"]
		}
		this.Kit.SetAngle(servoTilt, uint8(newTilt))

	case "Left":
		newPan := panAngle + cfg.PanTiltFactor
		if newPan > this.Kit.PanPos["left"] {
			newPan = this.Kit.PanPos["left"]
		}
		this.Kit.SetAngle(servoPan, uint8(newPan))

	case "Right":
		newPan := panAngle - cfg.PanTiltFactor
		if newPan < this.Kit.PanPos["right"] {
			newPan = this.Kit.PanPos["right"]
		}
		this.Kit.SetAngle(servoPan, uint8(newPan))

	case "CenterAll":
		this.Kit.SetCenter(servoPan)
		this.Kit.SetAngle(servoTilt, uint8(this.Kit.TiltPos["horizon"]))
	}
}
