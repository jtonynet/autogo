package main

import (
	"log"
	"net"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
	"gobot.io/x/gobot/platforms/raspi"

	application "github.com/jtonynet/autogo/application"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

//TODO env vars on viper
const (
	VERSION         = "v0.0.3"
	SERVOKIT_BUS    = 0
	SERVOKIT_ADDR   = 0x40
	ARDUINO_BUS     = 1
	ARDUINO_ADDR    = 0x18
	LCD_BUS         = 2
	LCD_ADDR        = 0x27
	LCD_COLLUMNS    = 16
	PAN_TILT_FACTOR = 30
)

func main() {
	r := raspi.NewAdaptor()
	keys := keyboard.NewDriver()

	///MOTORS
	motors := output.NewMotors(r)

	///SERVOKIT
	servoKit := output.NewServos(r, SERVOKIT_BUS, SERVOKIT_ADDR)
	servoPan := servoKit.Add("0", "pan")
	servoTilt := servoKit.Add("1", "tilt")

	///ARDUINO SONAR SET
	sonarSet, err := input.NewSonarSet(r, ARDUINO_BUS, ARDUINO_ADDR)
	if err != nil {
		log.Fatal(err)
	}

	///LCD
	lcd, err := output.NewLcd(LCD_BUS, LCD_ADDR, LCD_COLLUMNS)
	if err != nil {
		log.Fatal(err)
	}
	defer lcd.DeferAction()

	ip := GetOutboundIP()
	err = lcd.ShowMessage(string(ip), output.LINE_1)
	if err != nil {
		log.Fatal(err)
	}

	err = lcd.ShowMessage(VERSION+" Arrow key", output.LINE_2)
	if err != nil {
		log.Fatal(err)
	}

	work := func() {
		application.InitKeyboard(keys, motors, servoKit, sonarSet, lcd)
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
