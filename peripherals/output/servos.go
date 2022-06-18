package peripherals

import (
	"github.com/jtonynet/autogo/config"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

var TiltPos = map[string]int{
	"top":     0,
	"horizon": 130,
	"down":    180,
}

var PanPos = map[string]int{
	"left":  180,
	"right": 0,
}

type Servos struct {
	Driver  *i2c.PCA9685Driver
	Cfg     config.ServoKit
	kit     map[string]*gpio.ServoDriver
	TiltPos map[string]int
	PanPos  map[string]int
	PWMFreq float32
}

func NewServos(a *raspi.Adaptor, cfg config.ServoKit) *Servos {
	bus := cfg.Bus
	addr := cfg.Addr
	PWMFreq := cfg.PWMFrequency

	driver := i2c.NewPCA9685Driver(a,
		i2c.WithBus(bus),
		i2c.WithAddress(addr))

	kit := map[string]*gpio.ServoDriver{}
	this := &Servos{Driver: driver, Cfg: cfg, kit: kit, TiltPos: TiltPos, PanPos: PanPos, PWMFreq: PWMFreq}

	return this
}

func (this *Servos) Init() {
	this.Driver.SetPWMFreq(this.PWMFreq)
}

func (this *Servos) Add(servoId string, servoName string) *gpio.ServoDriver {
	s := gpio.NewServoDriver(this.Driver, servoId)
	s.SetName(servoName)

	this.kit[servoName] = s

	return s
}

func (this *Servos) GetByName(servoName string) *gpio.ServoDriver {
	return this.kit[servoName]
}

func (this *Servos) SetAngle(s *gpio.ServoDriver, angle uint8) {
	s.Move(angle)
}

func (this *Servos) SetCenter(s *gpio.ServoDriver) {
	s.Center()
}

func (this *Servos) SetMin(s *gpio.ServoDriver) {
	s.Min()
}

func (this *Servos) SetMax(s *gpio.ServoDriver) {
	s.Max()
}
