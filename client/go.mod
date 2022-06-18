module client

go 1.16

require (
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/gorilla/websocket v1.5.0
	github.com/jtonynet/autogo/config v0.0.0
	github.com/jtonynet/autogo/infrastructure v0.0.0
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/jtonynet/autogo/config => ../config

replace github.com/jtonynet/autogo/infrastructure => ../infrastructure
