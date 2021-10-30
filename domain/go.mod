module autogo

go 1.16

require (
	github.com/jtonynet/autogo/config v0.0.0
	github.com/jtonynet/autogo/peripherals v0.0.0
	gobot.io/x/gobot v1.15.0
)

replace github.com/jtonynet/autogo/config => ../config

replace github.com/jtonynet/autogo/peripherals => ../peripherals
