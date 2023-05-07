package domain

import (
	"fmt"
	"time"

	StatusDomain "github.com/jtonynet/autogo/domain/status"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	sensors "github.com/jtonynet/autogo/peripherals/sensors"
)

type IMU struct {
	IMU           *sensors.IMU
	MessageBroker *infrastructure.MessageBroker
	Status        *StatusDomain.Status
	Topic         string
	Delay         time.Duration
}

func NewIMU(imu *sensors.IMU, messageBroker *infrastructure.MessageBroker, status *StatusDomain.Status, topic string) *IMU {
	delay, _ := time.ParseDuration(imu.Cfg.Delay)
	this := &IMU{
		IMU:           imu,
		MessageBroker: messageBroker,
		Status:        status,
		Topic:         topic,
		Delay:         delay,
	}

	imu.Init()
	time.Sleep(time.Second * 5)

	return this
}

func (this *IMU) Worker() {
	for true {

		this.IMU.GetData()

		fmt.Println("Model", this.IMU.GetModel())
		fmt.Println("Accelerometer", this.IMU.GetAccelerometer())
		fmt.Println("Gyroscope", this.IMU.GetGyroscope())
		fmt.Println("Temperature", this.IMU.GetTemperature())

		fmt.Println()
		fmt.Println()

		time.Sleep(this.Delay)
	}
}
