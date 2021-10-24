package config

import "github.com/spf13/viper"

type Config struct {
	Version string `mapstructure:"VERSION"`

	ServokitBus          int     `mapstructure:"SERVOKIT_BUS"`
	ServokitAddr         int     `mapstructure:"SERVOKIT_ADDR"`
	ServokitPWMFrequency float32 `mapstructure:"SERVOKIT_PWM_FREQUENCY"`
	PanTiltFactor        int     `mapstructure:"PAN_TILT_FACTOR"`

	ArduinoBus        int     `mapstructure:"ARDUINO_BUS"`
	ArduinoAddr       int     `mapstructure:"ARDUINO_ADDR"`
	MinStopSonarValue float64 `mapstructure:"MIN_STOP_SONAR_VALUE"`

	LCDBus      int   `mapstructure:"LCD_BUS"`
	LCDAddr     uint8 `mapstructure:"LCD_ADDR"`
	LCDCollumns int   `mapstructure:"LCD_COLLUMNS"`

	MotorAPWMPin   string `mapstructure:"MOTOR_A_PWM_PIN"`
	MotorADir1Pin  string `mapstructure:"MOTOR_A_DIR1_PIN"`
	MotorADir2Pin  string `mapstructure:"MOTOR_A_DIR2_PIN"`
	MotorBPWMPin   string `mapstructure:"MOTOR_B_PWM_PIN "`
	MotorBDir1Pin  string `mapstructure:"MOTOR_B_DIR1_PIN"`
	MotorBDir2Pin  string `mapstructure:"MOTOR_B_DIR2_PIN"`
	MaxMotorsSpeed byte   `mapstructure:"MAX_MOTORS_SPEED"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
