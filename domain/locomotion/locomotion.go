package domain

import (
	LcdDomain "github.com/jtonynet/autogo/domain/lcd"
	StatusDomain "github.com/jtonynet/autogo/domain/status"
	actuators "github.com/jtonynet/autogo/peripherals/actuators"
)

type Locomotion struct {
	Motors *actuators.Motors
	Status *StatusDomain.Status
	LCD    *LcdDomain.LCD
}

func NewLocomotion(motors *actuators.Motors, lcd *LcdDomain.LCD, status *StatusDomain.Status) *Locomotion {
	this := &Locomotion{Motors: motors, LCD: lcd, Status: status}
	return this
}

func (this *Locomotion) Forward(speed byte) {
	this.Motors.Forward(speed)
}

func (this *Locomotion) Backward(speed byte) {
	this.Motors.Backward(speed)
}

func (this *Locomotion) Right(speed byte) {
	this.Motors.Right(speed)
}

func (this *Locomotion) Left(speed byte) {
	this.Motors.Left(speed)
}

func (this *Locomotion) Stop() {
	this.Motors.Stop()
}

func (this *Locomotion) Move(direction string) {
	oldDirection := this.Status.Direction
	cfg := this.Motors.Cfg

	switch direction {
	case "Front":
		if !this.Status.ColissionDetected {
			this.Forward(cfg.MaxSpeed)
			this.Status.Direction = "Front"
			this.Status.LCDMsg = this.Status.Direction
		}

	case "Back":
		this.Backward(cfg.MaxSpeed)
		this.Status.Direction = "Back"
		this.Status.LCDMsg = this.Status.Direction

	case "Right":
		this.Left(cfg.MaxSpeed)
		this.Status.Direction = "Right"
		this.Status.LCDMsg = this.Status.Direction

	case "Left":
		this.Right(cfg.MaxSpeed)
		this.Status.Direction = "Left"
		this.Status.LCDMsg = this.Status.Direction

	case "Stop":
		this.Stop()

		if this.Status.SonarSelfControll {
			this.Status.SonarSelfControll = false
			this.Status.SonarPreventCollision = true
		}

		this.Status.Direction = "Stop"
		this.Status.LCDMsg = this.Status.Version + " Arrow key"

	case "sonarPilot":
		if !this.Status.SonarSelfControll {
			//this.Stop()
			this.Status.SonarSelfControll = true
			this.Status.SonarPreventCollision = false
			this.Status.LCDMsg = this.Status.Version + " Snr Pilot"

		} else {
			this.Status.SonarSelfControll = false
			this.Status.SonarPreventCollision = true
			this.Status.LCDMsg = this.Status.Version + " Arrow key"
		}

	}

	if this.LCD != nil && oldDirection != this.Status.Direction {
		this.LCD.ShowMessage(this.Status.LCDMsg, 2)
	}
}
