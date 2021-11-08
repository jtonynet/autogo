package config

import "github.com/spf13/viper"

type ServoKit struct {
	Enabled       bool    `mapstructure:"SERVOKIT_ENABLED"`
	Bus           int     `mapstructure:"SERVOKIT_BUS"`
	Addr          int     `mapstructure:"SERVOKIT_ADDR"`
	PWMFrequency  float32 `mapstructure:"SERVOKIT_PWM_FREQUENCY"`
	PanTiltFactor int     `mapstructure:"PAN_TILT_FACTOR"`
}

type ArduinoSonar struct {
	Enabled      bool    `mapstructure:"ARDUINO_SONAR_ENABLED"`
	Bus          int     `mapstructure:"ARDUINO_SONAR_BUS"`
	Addr         int     `mapstructure:"ARDUINO_SONAR_ADDR"`
	MinStopValue float64 `mapstructure:"ARDUINO_MIN_SONAR_STOP_VALUE"`
}

type LCD struct {
	Enabled  bool  `mapstructure:"LCD_ENABLED"`
	Bus      int   `mapstructure:"LCD_BUS"`
	Addr     uint8 `mapstructure:"LCD_ADDR"`
	Collumns int   `mapstructure:"LCD_COLLUMNS"`
}

type Motors struct {
	Enabled  bool   `mapstructure:"MOTORS_ENABLED"`
	APWMPin  string `mapstructure:"MOTOR_A_PWM_PIN"`
	ADir1Pin string `mapstructure:"MOTOR_A_DIR1_PIN"`
	ADir2Pin string `mapstructure:"MOTOR_A_DIR2_PIN"`
	BPWMPin  string `mapstructure:"MOTOR_B_PWM_PIN"`
	BDir1Pin string `mapstructure:"MOTOR_B_DIR1_PIN"`
	BDir2Pin string `mapstructure:"MOTOR_B_DIR2_PIN"`
	MaxSpeed byte   `mapstructure:"MAX_MOTORS_SPEED"`
}

type Camera struct {
	Enabled bool   `mapstructure:"CAMERA_ENABLED"`
	Host    string `mapstructure:"CAMERA_STREAM_HOST"`
	Port    string `mapstructure:"CAMERA_STREAM_PORT"`
	Width   int    `mapstructure:"CAMERA_STREAM_WIDTH"`
}

type MQTT struct {
	Enabled  bool   `mapstructure:"MQTT_ENABLED"`
	Host     string `mapstructure:"MQTT_HOST "`
	Port     string `mapstructure:"MQTT_PORT "`
	User     string `mapstructure:"MQTT_USER "`
	Password string `mapstructure:"MQTT_PASSWORD"`
}

type Config struct {
	Version   string `mapstructure:"VERSION"`
	RobotName string `mapstructure:"ROBOT_NAME"`

	ServoKit     ServoKit     `mapstructure:",squash"`
	ArduinoSonar ArduinoSonar `mapstructure:",squash"`
	Camera       Camera       `mapstructure:",squash"`
	Motors       Motors       `mapstructure:",squash"`
	LCD          LCD          `mapstructure:",squash"`
	MQTT         MQTT         `mapstructure:",squash"`
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
