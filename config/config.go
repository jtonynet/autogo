package config

import "github.com/spf13/viper"

type ServoKit struct {
	Bus           int     `mapstructure:"SERVOKIT_BUS"`
	Addr          int     `mapstructure:"SERVOKIT_ADDR"`
	PWMFrequency  float32 `mapstructure:"SERVOKIT_PWM_FREQUENCY"`
	PanTiltFactor int     `mapstructure:"PAN_TILT_FACTOR"`
}

type ArduinoSonar struct {
	Bus          int     `mapstructure:"ARDUINO_BUS"`
	Addr         int     `mapstructure:"ARDUINO_ADDR"`
	MinStopValue float64 `mapstructure:"MIN_STOP_SONAR_VALUE"`
}

type LCD struct {
	Bus      int   `mapstructure:"LCD_BUS"`
	Addr     uint8 `mapstructure:"LCD_ADDR"`
	Collumns int   `mapstructure:"LCD_COLLUMNS"`
}

type Motors struct {
	APWMPin  string `mapstructure:"MOTOR_A_PWM_PIN"`
	ADir1Pin string `mapstructure:"MOTOR_A_DIR1_PIN"`
	ADir2Pin string `mapstructure:"MOTOR_A_DIR2_PIN"`
	BPWMPin  string `mapstructure:"MOTOR_B_PWM_PIN "`
	BDir1Pin string `mapstructure:"MOTOR_B_DIR1_PIN"`
	BDir2Pin string `mapstructure:"MOTOR_B_DIR2_PIN"`
	MaxSpeed byte   `mapstructure:"MAX_MOTORS_SPEED"`
}

type Config struct {
	Version string `mapstructure:"VERSION"`

	ServoKit     ServoKit     `mapstructure:",squash"`
	ArduinoSonar ArduinoSonar `mapstructure:",squash"`
	Motors       Motors       `mapstructure:",squash"`
	LCD          LCD          `mapstructure:",squash"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	return &cfg, nil
}
