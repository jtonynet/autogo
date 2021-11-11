package infrastructure

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	config "github.com/jtonynet/autogo/config"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//fmt.Printf("Received message: %s from topic %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost %v", err)
}

type MessageBroker struct {
	Client mqtt.Client
	Cfg    config.MessageBroker
}

func NewMessageBroker(cfg config.MessageBroker) *MessageBroker {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", cfg.Host, cfg.Port))
	opts.SetClientID("go_mqtt_client")

	if len(cfg.User) > 3 && len(cfg.Password) > 3 {
		opts.SetUsername(cfg.User)
		opts.SetPassword(cfg.Password)
	}

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	this := &MessageBroker{Client: client, Cfg: cfg}
	return this
}

func (this *MessageBroker) Disconnect() {
	this.Client.Disconnect(this.Cfg.WaitTTLDisconnect)
}

func (this *MessageBroker) Pub(topic string, message string) {
	token := this.Client.Publish(topic, 0, false, message)
	token.Wait()
}

func (this *MessageBroker) Sub(topic string) {
	token := this.Client.Subscribe(topic, 1, nil)
	token.Wait()
	//fmt.Printf("Subscribed to topic: %s ", topic)
}
