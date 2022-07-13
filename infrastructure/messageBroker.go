package infrastructure

import (
	"fmt"
	"time"

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
	opts.AddBroker(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))

	if len(cfg.User) > 3 && len(cfg.Password) > 3 && len(cfg.ClientID) > 3 {
		opts.SetUsername(cfg.User)
		opts.SetPassword(cfg.Password)
		opts.SetClientID(cfg.ClientID)
	}

	opts.SetDefaultPublishHandler(messagePubHandler)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {

		//panic(token.Error())
		fmt.Println(token.Error())
		fmt.Println("Retrying to connect in 1 second")
		time.Sleep(time.Second)

		return NewMessageBroker(cfg)
	}

	time.Sleep(time.Second)
	mb := &MessageBroker{Client: client, Cfg: cfg}

	return mb
}

func (mb *MessageBroker) Disconnect() {
	mb.Client.Disconnect(mb.Cfg.WaitTTLDisconnect)
}

func (mb *MessageBroker) Pub(topic string, message string) {

	token := mb.Client.Publish(topic, 0, false, message)
	token.Wait()
	//t := token.Wait()

	//if t {
	//	fmt.Println("MUITO CARALEO")
	//} else {
	//	fmt.Println("SINTO EM INFORMAR QUE DEU RUIM")
	//}

	//fmt.Println("\n-----------")
	//fmt.Printf("publish TEST %s on topic: %s ", message, topic)
	//fmt.Println("\n-----------")
}

func (mb *MessageBroker) Sub(topic string, receiverHandler func(mqtt.Client, mqtt.Message)) {
	if receiverHandler == nil {
		receiverHandler = defaultReceiver
	}

	token := mb.Client.Subscribe(topic, 1, receiverHandler)
	token.Wait()
	//fmt.Println("\n-----------")
	//fmt.Printf("Subscribed to topic: %s ", topic)
	//fmt.Println("\n-----------")
}

func defaultReceiver(client mqtt.Client, msg mqtt.Message) {
	msg.Ack()
	//output0 := "Robot.Control(\"default\" \"" + string(msg.Payload()) + "\")"
	//actuators := "message id:" + strconv.Itoa(int(msg.MessageID())) + " message = " + string(msg.Payload())
	//fmt.Println("\n++++++++++++++++")
	//fmt.Println(output0)
	//fmt.Println(actuators)
	//fmt.Println("\n++++++++++++++++")
}
