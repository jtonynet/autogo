package application

import (
	config "github.com/jtonynet/autogo/config"
	domain "github.com/jtonynet/autogo/domain"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

var (
	direction         string = ""
	lcdMsg            string = ""
	colissionDetected bool   = false
)

func Init(messageBroker *infrastructure.MessageBroker, kbd *input.Keyboard, motors *output.Motors, servoKit *output.Servos, lcd *output.Display, sonarSet *input.SonarSet, cfg *config.Config) {
	keys := kbd.Driver
	robotDomain := domain.NewRobot(messageBroker, motors, servoKit, lcd, sonarSet, cfg)

	keys.On(kbd.Key, func(data interface{}) {
		robotDomain.ControllByKeyboard(data)
	})
}
