package main

import (
	"fmt"
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

	var (
		botDevices []gobot.Device
		motors     *output.Motors  = nil
		servoKit   *output.Servos  = nil
		lcd        *output.Display = nil
		sonarSet   *input.SonarSet = nil
	)

	r := raspi.NewAdaptor()

	keys := keyboard.NewDriver()
	addDevice(&botDevices, keys)

	///MOTORS
	if cfg.Motors.Enabled {
		motors = output.NewMotors(r, cfg.Motors)
		addDevice(&botDevices, motors.MotorA)
		addDevice(&botDevices, motors.MotorB)
	}

	///SERVOKIT
	if cfg.ServoKit.Enabled {
		servoKit = output.NewServos(r, cfg.ServoKit)
		servoKit.Add("0", "pan")
		servoKit.Add("1", "tilt")

		addDevice(&botDevices, servoKit.Driver)
		addDevice(&botDevices, servoKit.GetByName("pan"))
		addDevice(&botDevices, servoKit.GetByName("tilt"))
	}

	fmt.Println(botDevices)

	///LCD
	if cfg.LCD.Enabled {
		lcd, err = output.NewLcd(cfg.LCD)
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
	}

	///ARDUINO SONAR SET
	if cfg.ArduinoSonar.Enabled {
		sonarSet, err = input.NewSonarSet(r, cfg.ArduinoSonar)
		if err != nil {
			log.Fatal(err)
		}
	}

	work := func() {
		application.Init(keys, motors, servoKit, lcd, sonarSet, cfg)
	}

	robot := gobot.NewRobot(
		cfg.RobotName,
		[]gobot.Connection{r},
		botDevices,
		work,
	)

	robot.Start()
}

func addDevice(deviceList *[]gobot.Device, device gobot.Device) {
	//Use only register gobot.Device
	*deviceList = append(*deviceList, device)
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
