package domain

type Status struct {
	ColissionDetected bool
	Direction         string
	LCDMsg            string
	MinStopValue      float64
	SonarData         map[string]float64
}
