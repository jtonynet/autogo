package domain

import (
	"fmt"
	"net"
	"strings"

	config "github.com/jtonynet/autogo/config"
	SonarDomain "github.com/jtonynet/autogo/domain/arduinoSonarSet"
	LcdDomain "github.com/jtonynet/autogo/domain/lcd"
	StatusDomain "github.com/jtonynet/autogo/domain/status"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"

	"gobot.io/x/gobot/platforms/keyboard"
)

type Robot struct {
	MessageBroker *infrastructure.MessageBroker

	Motors   *output.Motors
	ServoKit *output.Servos

	Cfg *config.Config

	LCD      *LcdDomain.LCD
	SonarSet *SonarDomain.Sonar
	Status   *StatusDomain.Status
}

func NewRobot(messageBroker *infrastructure.MessageBroker, motors *output.Motors, servoKit *output.Servos, display *output.Display, sonarSet *input.SonarSet, cfg *config.Config) *Robot {
	Status := &StatusDomain.Status{ColissionDetected: false, Direction: "", MinStopValue: cfg.ArduinoSonar.MinStopValue}
	this := &Robot{MessageBroker: messageBroker, Motors: motors, ServoKit: servoKit, Status: Status, Cfg: cfg}

	if servoKit != nil {
		servoPan := servoKit.GetByName("pan")
		servoTilt := servoKit.GetByName("tilt")

		servoKit.Init()
		servoKit.SetCenter(servoPan)
		servoKit.SetAngle(servoTilt, uint8(servoKit.TiltPos["horizon"]))
	}

	if display != nil {
		msgLine1 := getOutboundIP()
		if cfg.Camera.Enabled {
			s := []string{msgLine1, cfg.Camera.Port}
			msgLine1 = strings.Join(s, ":")
		}

		LCDTopic := fmt.Sprintf("%s/%s/lcd", cfg.ProjectName, cfg.RobotName)
		this.LCD = LcdDomain.NewLCD(display, messageBroker, LCDTopic)

		//TODO: Test only, remove after create robot client subscription
		if messageBroker != nil {
			messageBroker.Sub(LCDTopic)
		}

		this.LCD.ShowMessage(msgLine1, 1)
		this.LCD.ShowMessage(cfg.Version+" Arrow key", 2)
	}

	if sonarSet != nil {
		sonarTopic := fmt.Sprintf("%s/%s/sonar", cfg.ProjectName, cfg.RobotName)
		sonarDomain := SonarDomain.NewSonarSet(sonarSet, this.LCD, motors, messageBroker, Status, sonarTopic)
		this.SonarSet = sonarDomain

		//TODO: Test only, remove after create robot client subscription
		if messageBroker != nil {
			messageBroker.Sub(sonarTopic)
		}

		go sonarDomain.SonarWorker()
	}

	return this
}

func (this *Robot) ControllByKeyboard(data interface{}) {
	key := input.GetKeyEvent(data)

	if this.ServoKit != nil {
		go this.controllPanAndTilt(key.Key)
	}

	if this.Motors != nil {
		go this.controllMotors(key.Key)
	}
}

//TODO: Move to domain.ServoKit in future
func (this *Robot) controllPanAndTilt(k int) {
	cfg := this.Cfg
	servoPan := this.ServoKit.GetByName("pan")
	servoTilt := this.ServoKit.GetByName("tilt")

	panAngle := int(servoPan.CurrentAngle)
	tiltAngle := int(servoTilt.CurrentAngle)

	switch k {
	case keyboard.W:
		newTilt := tiltAngle - cfg.ServoKit.PanTiltFactor
		if newTilt < this.ServoKit.TiltPos["top"] {
			newTilt = this.ServoKit.TiltPos["top"]
		}
		this.ServoKit.SetAngle(servoTilt, uint8(newTilt))

	case keyboard.S:
		newTilt := tiltAngle + cfg.ServoKit.PanTiltFactor
		if newTilt > this.ServoKit.TiltPos["down"] {
			newTilt = this.ServoKit.TiltPos["down"]
		}
		this.ServoKit.SetAngle(servoTilt, uint8(newTilt))

	case keyboard.A:
		newPan := panAngle + cfg.ServoKit.PanTiltFactor
		if newPan > this.ServoKit.PanPos["left"] {
			newPan = this.ServoKit.PanPos["left"]
		}
		this.ServoKit.SetAngle(servoPan, uint8(newPan))

	case keyboard.D:
		newPan := panAngle - cfg.ServoKit.PanTiltFactor
		if newPan < this.ServoKit.PanPos["right"] {
			newPan = this.ServoKit.PanPos["right"]
		}
		this.ServoKit.SetAngle(servoPan, uint8(newPan))

	case keyboard.X:
		this.ServoKit.SetCenter(servoPan)
		this.ServoKit.SetAngle(servoTilt, uint8(this.ServoKit.TiltPos["horizon"]))
	}
}

//TODO: Move to domain.Motors in future
func (this *Robot) controllMotors(k int) {
	oldDirection := this.Status.Direction
	cfg := this.Cfg

	switch k {
	case keyboard.ArrowUp:
		if !this.Status.ColissionDetected {
			this.Motors.Forward(cfg.Motors.MaxSpeed)
			this.Status.Direction = "Front"
			this.Status.LCDMsg = this.Status.Direction
		}

	case keyboard.ArrowDown:
		this.Motors.Backward(cfg.Motors.MaxSpeed)
		this.Status.Direction = "Back"
		this.Status.LCDMsg = this.Status.Direction

	case keyboard.ArrowRight:
		this.Motors.Left(cfg.Motors.MaxSpeed)
		this.Status.Direction = "Right"
		this.Status.LCDMsg = this.Status.Direction

	case keyboard.ArrowLeft:
		this.Motors.Right(cfg.Motors.MaxSpeed)
		this.Status.Direction = "Left"
		this.Status.LCDMsg = this.Status.Direction

	case keyboard.Q:
		this.Motors.Stop()
		this.Status.Direction = ""
		this.Status.LCDMsg = cfg.Version + " Arrow key"
	}

	if this.LCD != nil && oldDirection != this.Status.Direction {
		this.LCD.ShowMessage(this.Status.LCDMsg, 2)
	}
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "offline"
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
