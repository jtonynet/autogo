package domain

import (
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	infrastructure "github.com/jtonynet/autogo/infrastructure"

	SonarDomain "github.com/jtonynet/autogo/domain/arduinoSonarSet"
	IMUDomain "github.com/jtonynet/autogo/domain/imu"
	LcdDomain "github.com/jtonynet/autogo/domain/lcd"
	LocomotionDomain "github.com/jtonynet/autogo/domain/locomotion"
	ServosDomain "github.com/jtonynet/autogo/domain/servos"
	StatusDomain "github.com/jtonynet/autogo/domain/status"

	config "github.com/jtonynet/autogo/config"

	actuators "github.com/jtonynet/autogo/peripherals/actuators"
	sensors "github.com/jtonynet/autogo/peripherals/sensors"
)

var (
	keyToRobotDirection = map[int]string{
		113: "Stop",
		65:  "Front",
		67:  "Right",
		66:  "Back",
		68:  "Left",
		112: "sonarPilot",
	}

	keyToCamDirection = map[int]string{
		119: "Top",
		100: "Right",
		115: "Down",
		97:  "Left",
		120: "CenterAll",
	}
)

type Robot struct {
	MessageBroker *infrastructure.MessageBroker

	LCD        *LcdDomain.LCD
	Locomotion *LocomotionDomain.Locomotion
	Servos     *ServosDomain.Servos
	SonarSet   *SonarDomain.Sonar
	IMU        *IMUDomain.IMU
	Status     *StatusDomain.Status

	Cfg *config.Config
}

func NewRobot(messageBroker *infrastructure.MessageBroker, motors *actuators.Motors, servos *actuators.Servos, display *actuators.Display, sonarSet *sensors.SonarSet, imu *sensors.IMU, cfg *config.Config) *Robot {
	Status := &StatusDomain.Status{
		ColissionDetected: false,
		Direction:         "Stop",
		Version:           cfg.Version,
		ProjectName:       cfg.ProjectName,
		RobotName:         cfg.RobotName,
		MinStopValue:      cfg.ArduinoSonar.MinStopValue,
	}

	this := &Robot{MessageBroker: messageBroker, Status: Status, Cfg: cfg}

	if servos != nil {
		servosDomain := ServosDomain.NewServos(servos)
		this.Servos = servosDomain
	}

	msgLine1 := infrastructure.GetOutboundIP()
	fmt.Println("\n" + msgLine1 + "\n")

	if display != nil {
		if cfg.Camera.Enabled {
			s := []string{msgLine1, cfg.Camera.Port}
			msgLine1 = strings.Join(s, ":")
		}

		LCDTopic := fmt.Sprintf("%s/%s/lcd", cfg.ProjectName, cfg.RobotName)
		this.LCD = LcdDomain.NewLCD(display, messageBroker, LCDTopic)

		//TODO: Test only, remove after create robot client subscription
		if messageBroker != nil {
			messageBroker.Sub(LCDTopic, nil)
		}

		this.LCD.ShowMessage(msgLine1, 1)
		this.LCD.ShowMessage(cfg.Version+" Arrow key", 2)
	}

	if motors != nil {
		locomotionDomain := LocomotionDomain.NewLocomotion(motors, this.LCD, this.Status)
		this.Locomotion = locomotionDomain
	}

	if sonarSet != nil {
		sonarTopic := fmt.Sprintf("%s/%s/sonar", cfg.ProjectName, cfg.RobotName)
		sonarDomain := SonarDomain.NewSonarSet(sonarSet, this.LCD, this.Locomotion, messageBroker, Status, sonarTopic)
		this.SonarSet = sonarDomain

		//TODO: Test only, remove after create robot client subscription
		if messageBroker != nil {
			go messageBroker.Sub(sonarTopic, nil)
		}

		this.Status.SonarSelfControll = false
		this.Status.SonarPreventCollision = true

		go sonarDomain.PreventCollisionWorker()
		//go sonarDomain.SelfControllWorker()
	}

	if imu != nil {
		imuTopic := fmt.Sprintf("%s/%s/imu", cfg.ProjectName, cfg.RobotName)
		imuDomain := IMUDomain.NewIMU(imu, messageBroker, Status, imuTopic)

		go imuDomain.Worker()
	}

	if messageBroker != nil {
		moveTopic := fmt.Sprintf("%s/%s/move", cfg.ProjectName, cfg.RobotName)
		go messageBroker.Sub(moveTopic, this.MoveMessageHandler)

		testTopic := fmt.Sprintf("%s/%s/test", cfg.ProjectName, cfg.RobotName)
		go messageBroker.Sub(testTopic, nil)

	}

	return this
}

func (this *Robot) ControllByKeyboard(data interface{}) {
	key := sensors.GetKeyEvent(data).Key

	var (
		action string
		exist  bool
	)

	if action, exist = keyToRobotDirection[key]; exist {
		this.controll("Direction", action)
	} else if action, exist = keyToCamDirection[key]; exist {
		this.controll("Cam", action)
	}
}

func (this *Robot) controll(command string, action string) {
	if command == "Direction" && this.Locomotion != nil {
		go this.Locomotion.Move(action)
	} else if command == "Cam" && this.Servos != nil {
		go this.Servos.ControllPanAndTilt(action)
	}
}

func (this *Robot) MoveMessageHandler(client mqtt.Client, msg mqtt.Message) {
	msg.Ack()

	output0 := "ROBOT:: this.Robot.Controll(\"Direction\", " + string(msg.Payload()) + "\")"
	fmt.Println(output0)
	this.controll("Direction", string(msg.Payload()))
}
