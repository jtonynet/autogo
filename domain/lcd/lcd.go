package domain

import (
	"encoding/json"
	"fmt"

	infrastructure "github.com/jtonynet/autogo/infrastructure"
	output "github.com/jtonynet/autogo/peripherals/output"
)

type LCD struct {
	Display       *output.Display
	MessageBroker *infrastructure.MessageBroker
	Topic         string
}

func NewLCD(display *output.Display, messageBroker *infrastructure.MessageBroker, topic string) *LCD {
	this := &LCD{Display: display, MessageBroker: messageBroker, Topic: topic}
	return this
}

func (this *LCD) ShowMessage(message string, line int) {
	this.Display.ShowMessage(message, line)

	if this.MessageBroker != nil {
		var (
			pubMsg = map[string]string{
				"message": message,
				"line":    fmt.Sprint(line),
			}
		)

		jsonMsg, err := json.Marshal(pubMsg)
		if err == nil {
			this.MessageBroker.Pub(this.Topic, string(jsonMsg))
		}

	}
}
