AutoGo
Raspberry Pi Golang Autonomous robot

Projeto em desenvolvimento com finalidade de estudo de uso de Golang com Robotica Básica e IOT (libs [Gobot](https://gobot.io/) e [periph.io](https://periph.io/)), funcional na plataforma base (Raspberry 3b+/4b Rasbpian recente).

![arquitetura](./docs/images/autogo_miro.png?raw=true "estrutura")

Em uma plataforma raspberry que siga a pinagem do esquema autoGo, instalar o pi-blaster
https://github.com/sarfata/pi-blaster

Rodar o comando para alterar pinagem i2c e habilitar todos pins para o gobot
'./autogo.sh'

Gerando Binario com raspberry como device alvo:
'GOARM=6 GOARCH=arm GOOS=linux go build main.go'

Rode o Binario, seja feliz

Com um teclado conectado ao raspberry, teclas de seta movimentam o veiculo, teclas "a, w, s, d" movimentam os servos da camera.
Stream de imagens, conducao pela web (pagina cliente com comunicacao via MQTT) e conducao autonoma (via sensor de sonar e outros metodos) estao sendo implementados

Pinagem do esquema autoGo
![esquema](./docs/images/autogo_fritzing_schema.jpg?raw=true "esquema")

![primeira versao](./docs/images/autogo_tank.jpg?raw=true "montado")

Referências:
  - https://gobot.io/documentation/platforms/raspi/
  - https://gobot.io/documentation/examples/firmata_motor/
  - https://pkg.go.dev/github.com/heupel/gobot/platforms/gpio#section-readme
  - https://github.com/d2r2/go-hd44780
  - https://github.com/hybridgroup/gobot/search?q=hd44780

-Próximas etapas:
  - Refatoração na estrutura do código
  - Condução Autônoma (Sonar set)
  - Condução por Fila MQTT (e web Socket)
  - Site Cliente para Condução
  - SH e makefile para automatizar dependencias em instalação nova
  - SH update de goversion no raspbian
  - Condução Autônoma (via Intel Neural Compute stick OU Google Coral)
  - Implantar ROS::: Golang
