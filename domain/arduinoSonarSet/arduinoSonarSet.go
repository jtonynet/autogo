package domain

import (
	"encoding/json"
	"fmt"
	"time"

	LcdDomain "github.com/jtonynet/autogo/domain/lcd"
	locomotionDomain "github.com/jtonynet/autogo/domain/locomotion"
	StatusDomain "github.com/jtonynet/autogo/domain/status"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	input "github.com/jtonynet/autogo/peripherals/input"
)

type Sonar struct {
	SonarSet      *input.SonarSet
	Locomotion    *locomotionDomain.Locomotion
	MessageBroker *infrastructure.MessageBroker
	Status        *StatusDomain.Status
	LCD           *LcdDomain.LCD
	Topic         string
	Delay         time.Duration
}

func NewSonarSet(sonarSet *input.SonarSet, LCD *LcdDomain.LCD, locomotion *locomotionDomain.Locomotion, messageBroker *infrastructure.MessageBroker, status *StatusDomain.Status, topic string) *Sonar {
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

func (this *Sonar) SelfControllWorker() {
	status := this.Status

	if this.Locomotion == nil {
		return
	}

	this.Locomotion.ControllMoviment("Front")

	for this.Status.SonarSelfControll {
		sonarData, err := this.SonarSet.GetData()
		if err != nil {
			this.SonarSet, _ = this.SonarSet.Reconnect()
			time.Sleep(this.Delay)
			go this.SelfControllWorker()
			return
		}

		if sonarData["center"] <= status.MinStopValue && status.Direction == "Front" && status.ColissionDetected == false {
			status.ColissionDetected = true

			this.Locomotion.Stop()

			if this.LCD != nil {
				s := fmt.Sprintf("STOP CRASH %.2f", sonarData["center"])
				this.LCD.ShowMessage(s, 2)
			}

			nextDirection := "Right"
			if sonarData["centerLeft"] > sonarData["centerRight"] {
				nextDirection = "Left"
			}

			status.ColissionDetected = false
			this.Locomotion.ControllMoviment(nextDirection)
			time.Sleep(700 * time.Millisecond)

		}

		if this.MessageBroker != nil {
			go this.sendDataToMessageBroker(sonarData)
		}

		status.SonarData = sonarData
		time.Sleep(this.Delay)
		this.Locomotion.ControllMoviment("Front")
	}
}

func (this *Sonar) PreventCrashWorker() {
	status := this.Status

	for true {
		sonarData, err := this.SonarSet.GetData()
		if err != nil {
			this.SonarSet, _ = this.SonarSet.Reconnect()
			time.Sleep(this.Delay)
			go this.PreventCrashWorker()
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
