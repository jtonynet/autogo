module peripherals

go 1.16

require (
	github.com/d2r2/go-hd44780 v0.0.0-20181002113701-74cc28c83a3e
	github.com/d2r2/go-i2c v0.0.0-20191123181816-73a8a799d6bc
	github.com/d2r2/go-logger v0.0.0-20210606094344-60e9d1233e22 // indirect
	github.com/hybridgroup/mjpeg v0.0.0-20140228234708-4680f319790e
	github.com/jtonynet/autogo/config v0.0.0-00010101000000-000000000000
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	gobot.io/x/gobot v1.15.0
	gocv.io/x/gocv v0.28.0
)

replace github.com/jtonynet/autogo/config => ../config
