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
	time.Sleep(1000 * time.Millisecond)
	if err != nil {
		return nil, err
	}

	this := &SonarSet{Master: a, conn: conn, Cfg: cfg}
	return this, nil
}

func (this *SonarSet) Reconnect() (sonarSet *SonarSet, err error) {
	fmt.Println("RECONNECT SONAR SET")
	return NewSonarSet(this.Master, this.Cfg)
}

func (this *SonarSet) GetData() (map[string]float64, error) {
	//fmt.Println("GET DATA A1")
	_, err := this.conn.Write([]byte("A"))
	if err != nil {
		return nil, err
	}

	//fmt.Println("GET DATA A2")
	sonarByteLen := 28
	buf := make([]byte, sonarByteLen)
	bytesRead, err := this.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	//fmt.Println("GET DATA A3")
	sonarData := ""
	if bytesRead == sonarByteLen {
		sonarData = string(buf[:])
	}

	//fmt.Println("GET DATA A4")
	dataValues := strings.Split(string(sonarData), ",")
	if len(dataValues) > 1 && len(datakeys) > len(dataValues)-1 {
		return nil, errors.New("sonar data dont match")
	}

	//fmt.Println("GET DATA A5")
	dataMap := make(map[string]float64)
	for i, data := range datakeys {
		dataMap[data], err = strconv.ParseFloat(dataValues[i], 64)
		if err != nil {
			//TODO: Customized error
			return nil, err
		}
	}

	//fmt.Println("GET DATA FINAL")
	fmt.Println(dataMap)
	return dataMap, nil
}
