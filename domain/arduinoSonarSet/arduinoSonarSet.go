package domain

import (
	"encoding/json"
	"fmt"
	"time"

	LcdDomain "github.com/jtonynet/autogo/domain/lcd"
	locomotionDomain "github.com/jtonynet/autogo/domain/locomotion"
	StatusDomain "github.com/jtonynet/autogo/domain/status"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	sensors "github.com/jtonynet/autogo/peripherals/sensors"
)

type Sonar struct {
	SonarSet      *sensors.SonarSet
	Locomotion    *locomotionDomain.Locomotion
	MessageBroker *infrastructure.MessageBroker
	Status        *StatusDomain.Status
	LCD           *LcdDomain.LCD
	Topic         string
	Delay         time.Duration
}

func NewSonarSet(sonarSet *sensors.SonarSet, LCD *LcdDomain.LCD, locomotion *locomotionDomain.Locomotion, messageBroker *infrastructure.MessageBroker, status *StatusDomain.Status, topic string) *Sonar {
	delay, _ := time.ParseDuration(sonarSet.Cfg.Delay)
	this := &Sonar{
		SonarSet:      sonarSet,
		LCD:           LCD,
		Locomotion:    locomotion,
		MessageBroker: messageBroker,
		Status:        status,
		Topic:         topic,
		Delay:         delay,
	}
	time.Sleep(time.Second * 5)
	return this
}

func (this *Sonar) sendDataToMessageBroker(sonarData map[string]float64) {
	j, err := json.Marshal(sonarData)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		this.MessageBroker.Pub(this.Topic, string(j))
	}
}

func (this *Sonar) Reconnect() {
	this.SonarSet, _ = this.SonarSet.Reconnect()
	time.Sleep(this.Delay)
	fmt.Println("Reconectei")
	return
}

func (this *Sonar) SelfControllWorker() {
	status := this.Status

	if this.Locomotion == nil {
		return
	}

	for true {
		if !this.Status.SonarSelfControll {
			time.Sleep(this.Delay)
			continue
		}

		this.Locomotion.Move("Front")

		sonarData, err := this.SonarSet.GetData()
		if err != nil {
			this.Reconnect()
			go this.SelfControllWorker()
			return
		}

		fmt.Println(sonarData)

		if sonarData["center"] <= status.MinStopValue && status.Direction == "Front" && status.ColissionDetected == false {
			status.ColissionDetected = true

			this.Locomotion.Stop()

			if this.LCD != nil {
				s := fmt.Sprintf("STOP CRASH %.2f", sonarData["center"])
				this.LCD.ShowMessage(s, 2)
			}

			this.Locomotion.Move("Back")
			time.Sleep(500 * time.Millisecond)

			sonarData, err = this.SonarSet.GetData()
			if err != nil {
				this.Locomotion.Move("Stop")
				this.Reconnect()
				go this.SelfControllWorker()
				return
			}

			nextDirection := "Right"
			if sonarData["centerLeft"] > sonarData["centerRight"] {
				nextDirection = "Left"
			}

			status.ColissionDetected = false
			this.Locomotion.Move(nextDirection)
			time.Sleep(700 * time.Millisecond)
			this.Locomotion.Move("Stop")
		}

		if this.MessageBroker != nil {
			go this.sendDataToMessageBroker(sonarData)
		}

		status.SonarData = sonarData
		time.Sleep(this.Delay)
		this.Locomotion.Move("Front")
	}
}

func (this *Sonar) PreventCollisionWorker() {
	status := this.Status

	for true {
		if !this.Status.SonarPreventCollision {
			time.Sleep(this.Delay)
			continue
		}

		sonarData, err := this.SonarSet.GetData()
		if err != nil {
			this.SonarSet, _ = this.SonarSet.Reconnect()
			time.Sleep(this.Delay)
			go this.PreventCollisionWorker()
			return
		}

		if sonarData["center"] <= status.MinStopValue && status.Direction == "Front" && status.ColissionDetected == false {
			status.ColissionDetected = true

			if this.Locomotion != nil {
				this.Locomotion.Stop()
			}

			if this.LCD != nil {
				s := fmt.Sprintf("STOP CRASH %.2f", sonarData["center"])
				this.LCD.ShowMessage(s, 2)
			}

		} else if status.ColissionDetected && status.Direction != "Front" {
			status.ColissionDetected = false
		}

		if this.MessageBroker != nil {
			go this.sendDataToMessageBroker(sonarData)
		}

		status.SonarData = sonarData
		time.Sleep(this.Delay)
	}
}
