cd ../pi-blaster/

sudo pi-blaster --gpio 17,18,16,19,13,20,12,21,22,5,23,4,24,25,26,27

#SERVOKIT PWM
sudo dtoverlay i2c-gpio bus=0 i2c_gpio_sda=12 i2c_gpio_scl=13 

#LCD
sudo dtoverlay i2c-gpio bus=2 i2c_gpio_sda=06 i2c_gpio_scl=07

#IMU
sudo dtoverlay i2c-gpio bus=3 i2c_gpio_sda=16 i2c_gpio_scl=20


echo "Initialize mosquitto mqtt/Websocket Server"
$(mosquitto -c /etc/mosquitto/mosquitto.conf)&
