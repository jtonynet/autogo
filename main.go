package main

import (
	"log"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"

	robotHandler "github.com/jtonynet/autogo/application"
	config "github.com/jtonynet/autogo/config"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	actuators "github.com/jtonynet/autogo/peripherals/actuators"
	sensors "github.com/jtonynet/autogo/peripherals/sensors"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	var (
		botDevices []gobot.Device

		motors   *actuators.Motors  = nil
		servoKit *actuators.Servos  = nil
		lcd      *actuators.Display = nil
		sonarSet *sensors.SonarSet  = nil
		imu      *sensors.IMU       = nil

		messageBroker *infrastructure.MessageBroker = nil
	)

	r := raspi.NewAdaptor()

	keys := sensors.GetKeyboard()
	addDevice(&botDevices, keys.Driver)

	if cfg.MessageBroker.Enabled {
		messageBroker = infrastructure.NewMessageBroker(cfg.MessageBroker)
	}

	///MOTORS
	if cfg.Motors.Enabled {
		motors = actuators.NewMotors(r, cfg.Motors)
		addDevice(&botDevices, motors.MotorA)
		addDevice(&botDevices, motors.MotorB)
	}

	///SERVOKIT
	if cfg.ServoKit.Enabled {
		servoKit = actuators.NewServos(r, cfg.ServoKit)
		servoKit.Add("0", "pan")
		servoKit.Add("1", "tilt")

		addDevice(&botDevices, servoKit.Driver)
		addDevice(&botDevices, servoKit.GetByName("pan"))
		addDevice(&botDevices, servoKit.GetByName("tilt"))
	}

	///LCD
	if cfg.LCD.Enabled {
		lcd, err = actuators.NewLcd(cfg.LCD)
		if err != nil {
			log.Fatal(err)
		}

		defer lcd.DeferAction()
	}

	///ARDUINO SONAR SET
	if cfg.ArduinoSonar.Enabled {
		sonarSet, err = sensors.NewSonarSet(r, cfg.ArduinoSonar)
		if err != nil {
			log.Fatal(err)
		}
	}

	///IMU
	if cfg.IMU.Enabled {
		imu = sensors.NewIMU(r, cfg.IMU)
		addDevice(&botDevices, imu.Driver)
	}

	work := func() {
		robotHandler.Init(messageBroker, keys, motors, servoKit, lcd, sonarSet, imu, cfg)

		///CAMERA STREAM
		if cfg.Camera.Enabled {
			go sensors.CameraServeStream(cfg.Camera)
		}
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
	*deviceList = append(*deviceList, device)
}
