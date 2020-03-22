package main

import (
	"github.com/swedishborgie/starlamp/lightctl"
	"time"
)

type SleepState int

const (
	StateUnknown SleepState = 0
	StateAwake   SleepState = 1
	StateAsleep  SleepState = 2
)

type State struct {
	AwakeTime      string              `json:"awakeTime"`
	AsleepTime     string              `json:"asleepTime"`
	NextAwakeTime  time.Time           `json:"nextAwakeTime"`
	NextAsleepTime time.Time           `json:"nextAsleepTime"`
	AwakeColor     lightctl.LightState `json:"awakeColor"`
	AsleepColor    lightctl.LightState `json:"asleepColor"`
	CurrentState   SleepState          `json:"currentState"`
	CurrentColor   lightctl.LightState `json:"currentColor"`
}

func (s SleepState) String() string {
	switch s {
	case StateAwake:
		return "awake"
	case StateAsleep:
		return "asleep"
	default:
		return "unknown"
	}
}
