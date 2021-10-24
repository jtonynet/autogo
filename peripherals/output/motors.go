package peripherals

import (
	"github.com/jtonynet/autogo/config"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

/*
// Objective: dual speed and direction control using MotorDriver
//
// | Enable | Dir 1 | Dir 2 | Motor         |
// +--------+-------+-------+---------------+
// | 0      | X     | X     | Off           |
// | 1      | 0     | 0     | 0ff           |
// | 1      | 0     | 1     | On (forward)  |
// | 1      | 1     | 0     | On (backward) |
// | 1      | 1     | 1     | Off           |

Motor Shield  | NodeMCU        | GPIO  | Purpose
--------------+----------------+-------+----------
A-Enable      | PWMA (Motor A) | 12	   | Speed
A-Dir1        | DIR1 (Motor A) | 15	   | Direction
A-Dir2        | DIR2 (Motor A) | 11	   | Direction
B-Enable      | PWMA (Motor B) | 35	   | Speed
B-Dir1        | DIR1 (Motor B) | 16	   | Direction
B-Dir2        | DIR2 (Motor B) | 18	   | Direction
*/

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
