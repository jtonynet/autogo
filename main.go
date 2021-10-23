package main

// Circuit: esp8266-and-l298n-motor-controller
// Objective: dual speed and direction control using MotorDriver
//
// | Enable | Dir 1 | Dir 2 | Motor         |
// +--------+-------+-------+---------------+
// | 0      | X     | X     | Off           |
// | 1      | 0     | 0     | 0ff           |
// | 1      | 0     | 1     | On (forward)  |
// | 1      | 1     | 0     | On (backward) |
// | 1      | 1     | 1     | Off           |

import (
	"fmt"
	"log"
	"net"
	"strings"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/keyboard"
	"gobot.io/x/gobot/platforms/raspi"

	deviceD2r2 "github.com/d2r2/go-hd44780"
	i2cD2r2 "github.com/d2r2/go-i2c"
)

/*
Motor Shield  | NodeMCU        | GPIO  | Purpose
--------------+----------------+-------+----------
A-Enable      | PWMA (Motor A) | 12	   | Speed
A-Dir1        | DIR1 (Motor A) | 15	   | Direction
A-Dir2        | DIR2 (Motor A) | 11	   | Direction
B-Enable      | PWMA (Motor B) | 35	   | Speed
B-Dir1        | DIR1 (Motor B) | 16	   | Direction
B-Dir2        | DIR2 (Motor B) | 18	   | Direction
*/

const (
	maPWMPin  = "12"
	maDir1Pin = "15"
	maDir2Pin = "11"
	mbPWMPin  = "35"
	mbDir1Pin = "16"
	mbDir2Pin = "18"
)

//TODO env vars on viper
const (
	VERSION      = "v0.0.4"
	LCD_COLLUMNS = 16
)

const (
	panTiltFactor = 30
)

const (
	maIndex = iota
	mbIndex
)

var (
	motorSpeed [2]byte
	motorInc   = [2]int{1, 1}
	counter    = [2]int{}
	motors     [2]*gpio.MotorDriver
)

func main() {
	r := raspi.NewAdaptor()
	keys := keyboard.NewDriver()

	///MOTORS
	motorA := gpio.NewMotorDriver(r, maPWMPin)
	motorA.ForwardPin = maDir1Pin
	motorA.BackwardPin = maDir2Pin
	motorA.SetName("Motor-A")

	motorB := gpio.NewMotorDriver(r, mbPWMPin)
	motorB.ForwardPin = mbDir1Pin
	motorB.BackwardPin = mbDir2Pin
	motorB.SetName("Motor-B")

	motors[maIndex] = motorA
	motors[mbIndex] = motorB
	///----

	///LCD
	//TODO: use lcd i2c gobot solution to 16x2 screen
	lcd, lcdI2cClose, err := lcdD2r2Factory()
	if err != nil {
		log.Fatal(err)
	}
	defer lcdI2cClose()

	err = lcd.BacklightOn()
	if err != nil {
		log.Fatal(err)
	}

	ip := GetOutboundIP()

	err = lcd.ShowMessage(string(ip), deviceD2r2.SHOW_LINE_1)
	if err != nil {
		log.Fatal(err)
	}

	err = lcd.ShowMessage(VERSION+" Arrow key", deviceD2r2.SHOW_LINE_2)
	if err != nil {
		log.Fatal(err)
	}
	///----

	///SERVOKIT
	servoDriver := i2c.NewPCA9685Driver(r,
		i2c.WithBus(0),
		i2c.WithAddress(0x40))

	pan := gpio.NewServoDriver(servoDriver, "0")
	tilt := gpio.NewServoDriver(servoDriver, "1")

	pan.SetName("pan")
	tilt.SetName("tilt")

	tiltPos := make(map[string]int)
	tiltPos["top"] = 0
	tiltPos["horizon"] = 130
	tiltPos["down"] = 180

	panPos := make(map[string]int)
	panPos["left"] = 180
	panPos["right"] = 0
	///----

	///ARDUINO SONAR SET
	arduinoConn, err := r.GetConnection(0x18, 1)
	if err != nil {
		log.Fatal(err)
	}
	///----

	firstRun := 1
	work := func() {
		servoDriver.SetPWMFreq(60)
		if firstRun == 1 {
			firstRun = 0
			pan.Center()
			tilt.Move(uint8(tiltPos["horizon"]))
		}

		keys.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			sonarData := ""
			if key.Key == keyboard.B {
				sonarData, err = getSonarData(arduinoConn)
				if err == nil {
					log.Println("///*********")
					log.Println("///Print arduino sonar data::")
					log.Println(sonarData)
					log.Println("///*********")
				}

			}

			panAngle := int(pan.CurrentAngle)
			tiltAngle := int(tilt.CurrentAngle)
			if key.Key == keyboard.W {
				newTilt := tiltAngle - panTiltFactor
				if newTilt < tiltPos["top"] {
					newTilt = tiltPos["top"]
				}
				tilt.Move(uint8(newTilt))

			} else if key.Key == keyboard.S {
				newTilt := tiltAngle + panTiltFactor
				if newTilt > tiltPos["down"] {
					newTilt = tiltPos["down"]
				}
				tilt.Move(uint8(newTilt))

			} else if key.Key == keyboard.A {
				newPan := panAngle + panTiltFactor
				if newPan > panPos["left"] {
					newPan = panPos["left"]
				}
				pan.Move(uint8(newPan))

			} else if key.Key == keyboard.D {
				newPan := panAngle - panTiltFactor
				if newPan < panPos["right"] {
					newPan = panPos["right"]
				}
				pan.Move(uint8(newPan))

			} else if key.Key == keyboard.X {
				pan.Center()
				tilt.Move(uint8(tiltPos["horizon"]))
			}

			if key.Key == keyboard.ArrowUp {
				motorA.Direction("forward")
				motorB.Direction("forward")
				motorA.Speed(255)
				motorB.Speed(255)
				lcd.ShowMessage(rightPad("Front", " ", LCD_COLLUMNS), deviceD2r2.SHOW_LINE_2)
			} else if key.Key == keyboard.ArrowDown {
				motorA.Direction("backward")
				motorB.Direction("backward")
				motorA.Speed(255)
				motorB.Speed(255)
				lcd.ShowMessage(rightPad("Back", " ", LCD_COLLUMNS), deviceD2r2.SHOW_LINE_2)
			} else if key.Key == keyboard.ArrowRight {
				motorA.Direction("forward")
				motorB.Direction("backward")
				motorA.Speed(255)
				motorB.Speed(255)
				lcd.ShowMessage(rightPad("Left", " ", LCD_COLLUMNS), deviceD2r2.SHOW_LINE_2)
			} else if key.Key == keyboard.ArrowLeft {
				motorA.Direction("backward")
				motorB.Direction("forward")
				motorA.Speed(255)
				motorB.Speed(255)
				lcd.ShowMessage(rightPad("Right", " ", LCD_COLLUMNS), deviceD2r2.SHOW_LINE_2)
			} else if key.Key == keyboard.Q {
				motorA.Speed(0)
				motorB.Speed(0)
				motorA.Direction("none")
				motorB.Direction("none")
				lcd.ShowMessage(VERSION+" Arrow key", deviceD2r2.SHOW_LINE_2)
			} else {
				fmt.Println("keyboard event!", key, key.Char)
			}
		})
	}

	robot := gobot.NewRobot(
		"my-robot",
		[]gobot.Connection{r},
		[]gobot.Device{
			motorA,
			motorB,
			keys,
			servoDriver,
			pan,
			tilt,
		},
		work,
	)

	robot.Start()
}

func lcdD2r2Factory() (*deviceD2r2.Lcd, func(), error) {
	// Create new connection to i2c-bus on 2 line with address 0x27.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2cD2r2.NewI2C(0x27, 2)
	if err != nil {
		log.Fatal(err)
	}

	// Construct lcd-device connected via I2C connection
	lcd, err := deviceD2r2.NewLcd(i2c, deviceD2r2.LCD_16x2)
	if err != nil {
		log.Fatal(err)
	}

	// Turn on the backlight
	err = lcd.BacklightOn()
	if err != nil {
		log.Fatal(err)
	}

	return lcd,
		func() {
			// Free I2C connection on exit
			defer i2c.Close()
		},
		nil
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

func rightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, (pLen-len(s)))
}

func getSonarData(sonarConn i2c.Connection) (string, error) {
	_, err := sonarConn.Write([]byte("A"))
	if err != nil {
		return "", err
	}

	sonarByteLen := 28
	buf := make([]byte, sonarByteLen)
	bytesRead, err := sonarConn.Read(buf)
	if err != nil {
		return "", err
	}

	sonarData := ""
	if bytesRead == sonarByteLen {
		sonarData = string(buf[:])
	}
	return sonarData, nil

}
