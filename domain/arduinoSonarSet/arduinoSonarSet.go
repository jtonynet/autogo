package domain

import (
	"encoding/json"
	"fmt"
	"time"

	LcdDomain "github.com/jtonynet/autogo/domain/lcd"
	StatusDomain "github.com/jtonynet/autogo/domain/status"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

type Sonar struct {
	SonarSet      *input.SonarSet
	Motors        *output.Motors //TODO: convert to domain
	MessageBroker *infrastructure.MessageBroker
	Status        *StatusDomain.Status
	LCD           *LcdDomain.LCD
	Topic         string
}

//TODO: Change output.Motors to domain.Motors in future
func NewSonarSet(SonarSet *input.SonarSet, LCD *LcdDomain.LCD, Motors *output.Motors, MessageBroker *infrastructure.MessageBroker, Status *StatusDomain.Status, Topic string) *Sonar {
	this := &Sonar{SonarSet: SonarSet, LCD: LCD, Motors: Motors, MessageBroker: MessageBroker, Status: Status, Topic: Topic}
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

func (this *Sonar) SonarWorker() {
	status := this.Status
	delayInMS, _ := time.ParseDuration(
		fmt.Sprintf("%vms", this.SonarSet.Cfg.DelayInMS))

	for true {
		sonarData, err := this.SonarSet.GetData()
		if err == nil {
			if sonarData["center"] <= status.MinStopValue && status.Direction == "Front" && status.ColissionDetected == false {
				status.ColissionDetected = true

				if this.Motors != nil {
					this.Motors.Stop()
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
			time.Sleep(delayInMS)
		}
	}
}
