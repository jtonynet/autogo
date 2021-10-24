package main

import (
	"log"
	"net"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
	"gobot.io/x/gobot/platforms/raspi"

	application "github.com/jtonynet/autogo/application"
	"github.com/jtonynet/autogo/config"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	r := raspi.NewAdaptor()
	keys := keyboard.NewDriver()

	///MOTORS
	motors := output.NewMotors(r, cfg.Motors)

	///SERVOKIT
	servoKit := output.NewServos(r, cfg.ServoKit)
	servoPan := servoKit.Add("0", "pan")
	servoTilt := servoKit.Add("1", "tilt")

	///ARDUINO SONAR SET
	sonarSet, err := input.NewSonarSet(r, cfg.ArduinoSonar)
	if err != nil {
		log.Fatal(err)
	}

	///LCD
	lcd, err := output.NewLcd(cfg.LCD)
	if err != nil {
		log.Fatal(err)
	}
	defer lcd.DeferAction()

	ip := GetOutboundIP()
	err = lcd.ShowMessage(string(ip), output.LINE_1)
	if err != nil {
		log.Fatal(err)
	}

	err = lcd.ShowMessage(cfg.Version+" Arrow key", output.LINE_2)
	if err != nil {
		log.Fatal(err)
	}

	work := func() {
		application.InitKeyboard(keys, motors, servoKit, sonarSet, lcd, cfg)
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
