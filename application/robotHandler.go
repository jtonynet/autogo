package application

import (
	config "github.com/jtonynet/autogo/config"
	aggregateDomain "github.com/jtonynet/autogo/domain"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	actuators "github.com/jtonynet/autogo/peripherals/actuators"
	sensors "github.com/jtonynet/autogo/peripherals/sensors"
)

var (
	direction         string = ""
	lcdMsg            string = ""
	colissionDetected bool   = false
)

func Init(messageBroker *infrastructure.MessageBroker, kbd *sensors.Keyboard, motors *actuators.Motors, servoKit *actuators.Servos, lcd *actuators.Display, sonarSet *sensors.SonarSet, imu *sensors.IMU, cfg *config.Config) {
	keys := kbd.Driver
	robotAggregate := aggregateDomain.NewRobot(messageBroker, motors, servoKit, lcd, sonarSet, imu, cfg)

	keys.On(kbd.Key, func(data interface{}) {
		robotAggregate.ControllByKeyboard(data)
	})
}
