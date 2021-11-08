###MQTT tbm deve ir para um SH... talvez esse faça sentido Dockerizar

- [Instalação MQTT](https://www.instructables.com/Installing-MQTT-BrokerMosquitto-on-Raspberry-Pi/)

- [Boas práticas MQTT Amazon](https://docs.aws.amazon.com/whitepapers/latest/designing-mqtt-topics-aws-iot-core/mqtt-design-best-practices.html)

- [Boas práticas MQTT Gerais](https://github.com/mqtt/mqtt.org/wiki/SYS-Topics)

- [Ferramenta MQTT](https://mqtt-explorer.com/)

- [Página de testes](https://test.mosquitto.org/)

- [Notas de instalação e resolução de problemas](https://github.com/eclipse/mosquitto/issues/529)

cd ~/etc/apt/sources.list.d
sudo -i

```
sharatkanthi commented on Apr 16, 2018
This is because the web socket library is not installed. Try running

sudo apt-get install libwebsockets-dev

before running

sudo apt-get install mosquitto
```

- [Problema de duplicidade, tenta subir o server MQTT duas vezes na porta 1883][http://www.steves-internet-guide.com/mosquitto-broker/]
```
mosquitto -V
sudo service mosquitto stop
sudo systemctl stop mosquitto.service
fuser -k 1883/tcp
mosquitto -V
```

- Testando (com 3 janelas distintas)
```
mosquitto

mosquitto_sub -h 0.0.0.0 -p 1883 -t 'autogo-test'

mosquitto_pub -h 0.0.0.0 -p 1883 -t 'autogo-test' -m 'test 3'
```

