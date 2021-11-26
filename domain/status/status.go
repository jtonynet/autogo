package domain

type Status struct {
	Version           string
	ProjectName       string
	RobotName         string
	ColissionDetected bool
	Direction         string
	LCDMsg            string
	MinStopValue      float64
	SonarData         map[string]float64
	SonarSelfControll bool
}
