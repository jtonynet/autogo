package peripherals

import (
	"github.com/jtonynet/autogo/config"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

var (
	motorSpeed [2]byte
	motorInc   = [2]int{1, 1}
	counter    = [2]int{}
	motors     [2]*gpio.MotorDriver
)

const (
	maIndex = iota
	mbIndex
)

type Motors struct {
	MotorA *gpio.MotorDriver
	MotorB *gpio.MotorDriver
}

func NewMotors(a *raspi.Adaptor, cfg config.Motors) *Motors {
	MotorA := gpio.NewMotorDriver(a, cfg.APWMPin)
	MotorA.ForwardPin = cfg.ADir1Pin
	MotorA.BackwardPin = cfg.ADir2Pin
	MotorA.SetName("Motor-A")

	MotorB := gpio.NewMotorDriver(a, cfg.BPWMPin)
	MotorB.ForwardPin = cfg.BDir1Pin
	MotorB.BackwardPin = cfg.BDir2Pin
	MotorB.SetName("Motor-B")

	this := &Motors{MotorA: MotorA, MotorB: MotorB}

	motors[maIndex] = MotorA
	motors[mbIndex] = MotorB

	return this
}

func (this *Motors) Forward(speed byte) {
	this.MotorA.Direction("forward")
	this.MotorB.Direction("forward")
	this.MotorA.Speed(speed)
	this.MotorB.Speed(speed)
}

func (this *Motors) Backward(speed byte) {
	this.MotorA.Direction("backward")
	this.MotorB.Direction("backward")
	this.MotorA.Speed(speed)
	this.MotorB.Speed(speed)
}

func (this *Motors) Right(speed byte) {
	this.MotorA.Direction("forward")
	this.MotorB.Direction("backward")
	this.MotorA.Speed(speed)
	this.MotorB.Speed(speed)
}

func (this *Motors) Left(speed byte) {
	this.MotorA.Direction("backward")
	this.MotorB.Direction("forward")
	this.MotorA.Speed(speed)
	this.MotorB.Speed(speed)
}

func (this *Motors) Stop() {
	this.MotorA.Speed(0)
	this.MotorB.Speed(0)
	this.MotorA.Direction("none")
	this.MotorB.Direction("none")
}
