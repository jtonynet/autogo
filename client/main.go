package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"github.com/jtonynet/autogo/config"
	"github.com/jtonynet/autogo/infrastructure"
)

var (
	addr              string
	topicBasePath     string
	messageBroker     *infrastructure.MessageBroker
	WebSocketConn     *websocket.Conn
	TopicsToSubscribe []string
)

type MessageToSend struct {
	Topic   string
	Message string
}

/**/
var upgrader = websocket.Upgrader{} // use default options

func proxyQueue(w http.ResponseWriter, r *http.Request) {
	webSocketConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer webSocketConn.Close()
	subscribeQueueTopics(webSocketConn)

	for {
		_, message, err := WebSocketConn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		m := string(message)
		var messageToSend MessageToSend
		json.Unmarshal([]byte(m), &messageToSend)

		topic := fmt.Sprintf("%s/%s", topicBasePath, messageToSend.Topic)
		messageBroker.Pub(topic, messageToSend.Message)

		time.Sleep(10 * time.Millisecond)
	}
}

func subscribeQueueTopics(ws *websocket.Conn) {
	if ws != nil {
		WebSocketConn = ws
		fmt.Println("WS connection exists")
	} else {
		fmt.Println("WS DONT connection exists")
		return
	}

	for _, topic := range TopicsToSubscribe {
		go messageBroker.Sub(topic, defaultReceiverHandler)
	}

}

func defaultReceiverHandler(client mqtt.Client, msg mqtt.Message) {
	if contains(TopicsToSubscribe, msg.Topic()) {
		m := fmt.Sprintf("{\"topic\":\"%s\", \"message\":%s}", msg.Topic(), string(msg.Payload()))
		fmt.Println(m)
		WebSocketConn.WriteMessage(websocket.TextMessage, []byte(m))
	}
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

/**/

func main() {

	cfg, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	var port string
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = string(cfg.Client.Port)
	}
	addr = fmt.Sprintf(":%s", port)

	/**/
	messageBroker = infrastructure.NewMessageBroker(cfg.MessageBroker)

	topicBasePath = fmt.Sprintf("%s/%s", cfg.ProjectName, cfg.RobotName)
	topics := strings.Split(cfg.Client.TopicsToSubscribe, ",")
	for _, topic := range topics {
		topicComplete := fmt.Sprintf("%s/%s", topicBasePath, topic)
		TopicsToSubscribe = append(TopicsToSubscribe, topicComplete)
	}

	//flag.Parse()
	//log.SetFlags(0)
	http.HandleFunc("/proxy_queue", proxyQueue)
	/**/
	static := http.FileServer(http.Dir("./static"))
	http.Handle("/", static)

	ip := infrastructure.GetOutboundIP()
	fmt.Println("OK on ", ip, addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
