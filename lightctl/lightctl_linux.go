package lightctl

// #cgo LDFLAGS: -lwiringPi
// #include <wiringPi.h>
// #include <unistd.h>
// #include <stdio.h>
// #define TOGGLE_PIN 1
// #define STATE_COUNT 9
// #define RESET_DELAY (10 * 1000)
// #define MIN_DELAY (100 * 1000)
//
// int currentState = 0;
//
// void init(void) {
//	wiringPiSetup();
//  pinMode(TOGGLE_PIN, OUTPUT);
// }
//
// void reset(void) {
//	for (int i = 0; i < 100; i++) {
//		digitalWrite(TOGGLE_PIN, LOW);
//		usleep(RESET_DELAY);
//		digitalWrite(TOGGLE_PIN, HIGH);
//		usleep(RESET_DELAY);
//	}
//  currentState = 0;
//	usleep(MIN_DELAY);
// }
//
// void next(void) {
// 	digitalWrite(TOGGLE_PIN, LOW);
//	usleep(MIN_DELAY);
//  digitalWrite(TOGGLE_PIN, HIGH);
//  currentState = (currentState + 1) % STATE_COUNT;
// }
//
// void setState(int state) {
// 	while (currentState != state) {
// 		next();
//		usleep(MIN_DELAY);
//	}
// }
//
// int getState(void) {
// 	return currentState;
// }
import "C"

func init() {
	C.init()
}

func Next() {
	C.next()
}

func Reset() {
	C.reset()
}

func SetState(state LightState) {
	C.setState(C.int(state))
}

func GetState() LightState {
	return LightState(C.getState())
}
