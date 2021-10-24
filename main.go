package main

import (
	"log"
	"net"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
	"gobot.io/x/gobot/platforms/raspi"

	application "github.com/jtonynet/autogo/application"
	config "github.com/jtonynet/autogo/config"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	r := raspi.NewAdaptor()
	keys := keyboard.NewDriver()

	///MOTORS
	motors := output.NewMotors(r, config)

	///SERVOKIT
	servoKit := output.NewServos(r, config.ServokitBus, config.ServokitAddr, config.ServokitPWMFrequency)
	servoPan := servoKit.Add("0", "pan")
	servoTilt := servoKit.Add("1", "tilt")

	///ARDUINO SONAR SET
	sonarSet, err := input.NewSonarSet(r, config.ArduinoBus, config.ArduinoAddr)
	if err != nil {
		log.Fatal(err)
	}

	///LCD
	lcd, err := output.NewLcd(config.LCDBus, config.LCDAddr, config.LCDCollumns)
	if err != nil {
		log.Fatal(err)
	}
	defer lcd.DeferAction()

	ip := GetOutboundIP()
	err = lcd.ShowMessage(string(ip), output.LINE_1)
	if err != nil {
		log.Fatal(err)
	}

	err = lcd.ShowMessage(config.Version+" Arrow key", output.LINE_2)
	if err != nil {
		log.Fatal(err)
	}

	work := func() {
		application.InitKeyboard(keys, motors, servoKit, sonarSet, lcd, config)
	}

	robot := gobot.NewRobot(
		"my-robot",
		[]gobot.Connection{r},
		[]gobot.Device{
			motors.MotorA,
			motors.MotorB,
			keys,
			servoKit.Driver,
			servoPan,
			servoTilt,
		},
		work,
	)

	robot.Start()
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "ip offline"
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
