package gamelogic

import "time"

type ClockState struct {
	Timestamp      time.Time
	BlackOTCount   uint8
	BlackTimerRead uint32
	WhiteOTCount   uint8
	WhiteTimerRead uint32
	TimeLimit      uint32
	OTCout         uint8
}
