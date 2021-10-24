//Programa: Coleta dados de um poll de sensores Ultrassonico 
//HC-SR04 e comunica via I2C Arduino e Raspberry Pi (permite 4 sensores)

#include <Wire.h>
#include <Ultrasonic.h>

//Define os pinos para o trigger e echo
#define pino_trigger_front 4
#define pino_echo_front 5

#define pino_trigger_front_right 2
#define pino_echo_front_right 3

#define pino_trigger_back 7
#define pino_echo_back 6

#define pino_trigger_front_left 9
#define pino_echo_front_left 8

//Inicializa o sensor nos pinos definidos acima
Ultrasonic ultrasonic_front(pino_trigger_front, pino_echo_front);
Ultrasonic ultrasonic_front_right(pino_trigger_front_right, pino_echo_front_right);
Ultrasonic ultrasonic_back(pino_trigger_back, pino_echo_back);
Ultrasonic ultrasonic_front_left(pino_trigger_front_left, pino_echo_front_left);

char resultBuffer[36];
char col_c[7];
char col_r[7];
char col_b[7];
char col_l[7];
char value_c[4];
char value_r[4];
char value_b[4];
char value_l[4];

void setup()
{
  Serial.begin(9600);
  Wire.begin(0x18);
  Wire.onRequest(requestEvent);
}

void requestEvent()
{
  String result = "";
  float max_limit = 999.99;

  //Its dont work, i use string concat fix solution
  //sprintf(str, "%s,%s,%s,%s,", col_c, col_r, col_b, col_l);
  
  float cmMsec_front;
  long microsec_front = ultrasonic_front.timing();
  cmMsec_front = ultrasonic_front.convert(microsec_front, Ultrasonic::CM);

  if(cmMsec_front > max_limit){
    cmMsec_front = max_limit;
  }
  
  dtostrf(cmMsec_front,4,2, value_c);
  sprintf(col_c, "%s,", value_c);
  result.concat(col_c);

  float cmMsec_front_right;
  long microsec_front_right = ultrasonic_front_right.timing();
  cmMsec_front_right = ultrasonic_front_right.convert(microsec_front_right, Ultrasonic::CM);

  if(cmMsec_front_right > max_limit){
    cmMsec_front_right = max_limit;
  }
  
  dtostrf(cmMsec_front_right,4,2, value_r);
  sprintf(col_r, "%s,", value_r);
  result.concat(col_r);

  float cmMsec_back;
  long microsec_back = ultrasonic_back.timing();
  cmMsec_back = ultrasonic_back.convert(microsec_back, Ultrasonic::CM);

  if(cmMsec_back > max_limit){
    cmMsec_back = max_limit;
  }
  
  dtostrf(cmMsec_back,4,2, value_b);
  sprintf(col_b, "%s,", value_b);
  result.concat(col_b);

  float cmMsec_front_left;
  long microsec_front_left = ultrasonic_front_left.timing();
  cmMsec_front_left = ultrasonic_front_left.convert(microsec_front_left, Ultrasonic::CM);

  if(cmMsec_front_left > max_limit){
    cmMsec_front_left = max_limit;
  }
  
  dtostrf(cmMsec_front_left,4,2, value_l);
  sprintf(col_l, "%s,", value_l);
  result.concat(col_l);
  
  //        center,cRight,back  ,cLeft 
  //result="000.00,000.00,000.00,000.00,";
  //result="999.99,999.99,999.99,999.99,";
  
  result.toCharArray(resultBuffer, 36);
  
  Serial.println(result);
  Serial.println(" ");
  Wire.write(resultBuffer);
}

void loop()
{
  delay(1000);
}