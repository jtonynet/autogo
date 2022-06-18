package peripherals

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jtonynet/autogo/config"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

var datakeys = []string{"center", "centerRight", "back", "centerLeft"}

type SonarSet struct {
	Master *raspi.Adaptor
	conn   i2c.Connection
	Cfg    config.ArduinoSonar
}

func NewSonarSet(a *raspi.Adaptor, cfg config.ArduinoSonar) (sonarSet *SonarSet, err error) {
	bus := cfg.Bus
	addr := cfg.Addr
	conn, err := a.GetConnection(addr, bus)
	time.Sleep(1 * time.Millisecond)
	if err != nil {
		return nil, err
	}

	this := &SonarSet{Master: a, conn: conn, Cfg: cfg}
	return this, nil
}

func (this *SonarSet) Reconnect() (sonarSet *SonarSet, err error) {
	return NewSonarSet(this.Master, this.Cfg)
}

func (this *SonarSet) GetData() (map[string]float64, error) {
	_, err := this.conn.Write([]byte("A"))
	if err != nil {
		return nil, err
	}

	sonarByteLen := 28
	buf := make([]byte, sonarByteLen)
	bytesRead, err := this.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	sonarData := ""
	if bytesRead == sonarByteLen {
		sonarData = string(buf[:])
	}

	fmt.Println(sonarData)

	dataValues := strings.Split(string(sonarData), ",")
	if len(dataValues) > 1 && len(datakeys) > len(dataValues)-1 {
		return nil, errors.New("sonar data dont match")
	}

	dataMap := make(map[string]float64)
	for i, data := range datakeys {
		dataMap[data], err = strconv.ParseFloat(dataValues[i], 64)
		if err != nil {
			//TODO: Customized error
			return nil, err
		}
	}

	return dataMap, nil
}
