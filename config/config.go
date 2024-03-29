package config

import (
	"github.com/spf13/viper"
	"os"
)

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
	Delay        string  `mapstructure:"ARDUINO_SONAR_DELAY"`
}

type LCD struct {
	Enabled  bool  `mapstructure:"LCD_ENABLED"`
	Bus      int   `mapstructure:"LCD_BUS"`
	Addr     uint8 `mapstructure:"LCD_ADDR"`
	Collumns int   `mapstructure:"LCD_COLLUMNS"`
}

type IMU struct {
	Enabled bool   `mapstructure:"IMU_ENABLED"`
	Model   string `mapstructure:"IMU_MODEL"`
	Bus     uint8  `mapstructure:"IMU_BUS"`
	Addr    uint8  `mapstructure:"IMU_ADDR"`
	Delay   string `mapstructure:"IMU_DELAY"`
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

type MessageBroker struct {
	Enabled           bool   `mapstructure:"MESSAGEBROKER_ENABLED"`
	Host              string `mapstructure:"MESSAGEBROKER_HOST"`
	Port              string `mapstructure:"MESSAGEBROKER_PORT"`
	ClientID          string `mapstructure:"MESSAGEBROKER_CLIENT_ID"`
	User              string `mapstructure:"MESSAGEBROKER_USER"`
	Password          string `mapstructure:"MESSAGEBROKER_PASSWORD"`
	WaitTTLDisconnect uint   `mapstructure:"MESSAGEBROKER_TTL_DISCONNECT_IN_MS"`
}

type Client struct {
	Port              string `mapstructure:"CLIENT_PORT"`
	TopicsToSubscribe string `mapstructure:"CLIENT_TOPICS_SUBSCRIBE"`
}

type Config struct {
	Version     string `mapstructure:"VERSION"`
	ProjectName string `mapstructure:"PROJECT_NAME"`
	RobotName   string `mapstructure:"ROBOT_NAME"`

	ServoKit      ServoKit      `mapstructure:",squash"`
	ArduinoSonar  ArduinoSonar  `mapstructure:",squash"`
	Camera        Camera        `mapstructure:",squash"`
	Motors        Motors        `mapstructure:",squash"`
	LCD           LCD           `mapstructure:",squash"`
	IMU           IMU           `mapstructure:",squash"`
	MessageBroker MessageBroker `mapstructure:",squash"`
	Client        Client        `mapstructure:",squash"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError, *os.PathError:
			// NOTE: Need to log out to console regardless of log level
			//log.Info("using config from env vars instead config file")
		default:
			return nil, err
		}
	}

	viper.AutomaticEnv()

	var cfg Config
	err = viper.Unmarshal(&cfg)
	return &cfg, nil
}
