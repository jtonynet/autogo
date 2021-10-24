package peripherals

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jtonynet/autogo/config"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

var datakeys = []string{"center", "centerRight", "back", "centerLeft"}

type SonarSet struct {
	conn i2c.Connection
}

func NewSonarSet(a *raspi.Adaptor, cfg config.ArduinoSonar) (sonarSet *SonarSet, err error) {
	bus := cfg.Bus
	addr := cfg.Addr
	conn, err := a.GetConnection(addr, bus)
	if err != nil {
		return nil, err
	}

	this := &SonarSet{conn: conn}
	return this, nil
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

	dataValues := strings.Split(string(sonarData), ",")
	if len(dataValues) > 1 && len(datakeys) > len(dataValues)-1 {
		return nil, errors.New("sonar data dont match")
	}

	dataMap := make(map[string]float64)
	for i, data := range datakeys {
		dataMap[data], err = strconv.ParseFloat(dataValues[i], 64)
		if err != nil {
			//TODO: customized error
			return nil, err
		}
	}

	return dataMap, nil
}
