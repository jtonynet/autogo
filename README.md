AutoGo
Raspberry Pi Golang Autonomous robot

Projeto em desenvolvimento, funcional na plataforma base (Raspberry 3b+/4b Rasbpian recente).

![arquitetura](./docs/images/autogo_miro.png?raw=true "estrutura")

Em uma plataforma raspberry que siga a pinagem do esquema autoGo, instalar o pi-blaster
https://github.com/sarfata/pi-blaster

Rodar o comando para alterar pinagem i2c e habilitar todos pins para o gobot
'./autogo.sh'

Gerando Binario com raspberry como device alvo:
'GOARM=6 GOARCH=arm GOOS=linux go build main.go'

Rode o Binario, seja feliz

Pinagem do esquema autoGo
![esquema](./docs/images/autogo_fritzing_schema.jpg?raw=true "esquema")

Versao Inicial do autoGo
![primeira versao](./docs/images/autogo_tank.jpg?raw=true "montado")

Referências:
  - https://gobot.io/documentation/platforms/raspi/
  - https://gobot.io/documentation/examples/firmata_motor/
  - https://pkg.go.dev/github.com/heupel/gobot/platforms/gpio#section-readme
  - https://github.com/d2r2/go-hd44780
  - https://github.com/hybridgroup/gobot/search?q=hd44780

i2c:
  - JHD1313M1 LCD Display w/RGB Backlight
  - https://github.com/hybridgroup/gobot/blob/a8f33b2fc012951104857c485e85b35bf5c4cb9d/drivers/i2c/README.md

-Próximas etapas (deixar identico ao ultimo Master da versão Python):
  - Refatoração na estrutura do código
  - Comunicação com arduino (Sonar set)
  - sh e makefile para automatizar dependencias em instalação nova
  - sh update de goversion no raspbian
