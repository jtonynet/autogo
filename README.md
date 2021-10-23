Autobot-go WIP

Log da versão
  - Motores funcionais
  - Respondendo ao teclad w,a,s,d (q para parar)
  - LCD funcional com lib d2r2 (descobrir como funciona e migrar para lib padrão gobot :/ )


Em uma plataforma raspberry que siga a pinagem do esquema autobot, instalar o pi-blaster
https://github.com/sarfata/pi-blaster

Rodar o comando para alterar pinagem i2c e habilitar todos pins para o gobot
'./autogo.sh'

Gerando Binario com raspberry como device alvo:
'GOARM=6 GOARCH=arm GOOS=linux go build main.go'

Rode o Binario, seja feliz

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