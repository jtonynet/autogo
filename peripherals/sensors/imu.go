package peripherals

import (
	"time"

	"github.com/jtonynet/autogo/config"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

type IMU struct {
	Driver *i2c.MPU6050Driver
	Cfg    config.IMU
	delay  time.Duration
}

func NewIMU(a *raspi.Adaptor, cfg config.IMU) *IMU {
	bus := cfg.Bus
	addr := cfg.Addr

	driver := i2c.NewMPU6050Driver(a,
		i2c.WithBus(int(bus)),
		i2c.WithAddress(int(addr)))

	delay, _ := time.ParseDuration(cfg.Delay)
	this := &IMU{Driver: driver, Cfg: cfg, delay: delay}

	return this
}

func (this *IMU) Init() {
	this.Driver.Start()
}

func (this *IMU) GetData() {
	this.Driver.GetData()
}

func (this *IMU) GetAccelerometer() i2c.ThreeDData {
	return this.Driver.Accelerometer
}

func (this *IMU) GetGyroscope() i2c.ThreeDData {
	return this.Driver.Gyroscope
}

func (this *IMU) GetTemperature() int16 {
	return this.Driver.Temperature
}

func (this *IMU) GetModel() string {
	return this.Cfg.Model
}
