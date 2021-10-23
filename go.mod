module autogo

go 1.16

require (
	github.com/matrixreality/autogo/handlers v0.0.0-00010101000000-000000000000
	github.com/matrixreality/autogo/peripherals v0.0.0
	gobot.io/x/gobot v1.15.0
)

replace github.com/matrixreality/autogo/peripherals => ./peripherals

replace github.com/matrixreality/autogo/handlers => ./handlers
