package main

import (
	"log"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"

	application "github.com/jtonynet/autogo/application"
	config "github.com/jtonynet/autogo/config"
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

	keys := input.GetKeyboard()
	addDevice(&botDevices, keys.Driver)

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

	///LCD
	if cfg.LCD.Enabled {
		lcd, err = output.NewLcd(cfg.LCD)
		if err != nil {
			log.Fatal(err)
		}
		defer lcd.DeferAction()
	}

	///ARDUINO SONAR SET
	if cfg.ArduinoSonar.Enabled {
		sonarSet, err = input.NewSonarSet(r, cfg.ArduinoSonar)
		if err != nil {
			log.Fatal(err)
		}
	}

	///CAMERA STREAM
	if cfg.Camera.Enabled {
		go input.CameraServeStream(cfg.Camera)
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
