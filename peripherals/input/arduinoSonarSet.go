package peripherals

import (
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

type SonarSet struct {
	conn i2c.Connection
}

func NewSonarSet(a *raspi.Adaptor, bus int, addr uint8) (sonarSet *SonarSet, err error) {
	conn, err := a.GetConnection(0x18, 1)
	if err != nil {
		return nil, err
	}

	this := &SonarSet{conn: conn}
	return this, nil
}

func (this *SonarSet) GetData() (string, error) {
	_, err := this.conn.Write([]byte("A"))
	if err != nil {
		return "", err
	}

	sonarByteLen := 28
	buf := make([]byte, sonarByteLen)
	bytesRead, err := this.conn.Read(buf)
	if err != nil {
		return "", err
	}

	sonarData := ""
	if bytesRead == sonarByteLen {
		sonarData = string(buf[:])
	}

	return sonarData, nil
}
