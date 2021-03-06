package gameLogic_test

import (
	"strings"
	"testing"
	"time"

	gl "github.com/alanzeng6181/game-of-go/gameLogic"
)

func TestClock_5_2_2_SameColor(t *testing.T) {
	input := make(chan gl.Stone)
	clock, err := gl.NewClock(gl.MakeClockConfig(5*time.Millisecond, 2*time.Millisecond, 2), input)
	if err != nil {
		t.Error(err)
	}
	go func(send chan<- gl.Stone) {
		time.Sleep(1 * time.Millisecond)
		send <- gl.Black
		time.Sleep(1 * time.Millisecond)
		send <- gl.Black
	}(input)

	clock.Start()
	timeout, err := clock.GetTimeout()
	if err != nil {
		t.Errorf("could not get clock's timeout channel => %v", err)
		return
	}

	correct := false
	select {
	case <-time.NewTimer(4 * time.Millisecond).C:
	case result := <-timeout:
		if strings.Contains(result.Error(), "expected different color") {
			correct = true
		}
	}
	if !correct {
		t.Error("expected sending same color twice to clock input causing error, but it didn't")
	}
}

func TestClock_5_2_2_BlackTimout(t *testing.T) {
	input := make(chan gl.Stone)
	clock, err := gl.NewClock(gl.MakeClockConfig(5*time.Millisecond, 2*time.Millisecond, 2), input)
	if err != nil {
		t.Error(err)
	}
	go func(send chan<- gl.Stone) {
		time.Sleep(1 * time.Millisecond)
		send <- gl.Black
		time.Sleep(1 * time.Millisecond)
		send <- gl.White
	}(input)

	clock.Start()
	timeout, err := clock.GetTimeout()
	if err != nil {
		t.Errorf("could not get clock's timeout channel => %v", err)
		return
	}

	correct := false
	select {
	case <-time.NewTimer(20 * time.Millisecond).C:
	case result := <-timeout:
		if strings.Contains(result.Error(), "timeout") && result.Stone == gl.Black {
			correct = true
		}
	}
	if !correct {
		t.Error("expected black timeout, but it didn't")
	}
}

func TestClock_5_2_2_NoTimout(t *testing.T) {
	input := make(chan gl.Stone)
	clock, err := gl.NewClock(gl.MakeClockConfig(5*time.Millisecond, 2*time.Millisecond, 2), input)
	if err != nil {
		t.Error(err)
	}
	go func(send chan<- gl.Stone) {
		time.Sleep(5 * time.Millisecond)
		send <- gl.Black
		send <- gl.White
		time.Sleep(3 * time.Millisecond)
		send <- gl.Black
		send <- gl.White
		time.Sleep(3 * time.Millisecond)
		send <- gl.Black
		send <- gl.White
		time.Sleep(1 * time.Millisecond)
		send <- gl.Black
		send <- gl.White
		time.Sleep(1 * time.Millisecond)
		send <- gl.Black
		send <- gl.White
		time.Sleep(1 * time.Millisecond)
		send <- gl.Black
		send <- gl.White
		time.Sleep(1 * time.Millisecond)
		send <- gl.Black
		send <- gl.White
	}(input)

	clock.Start()
	timeout, err := clock.GetTimeout()
	if err != nil {
		t.Errorf("could not get clock's timeout channel => %v", err)
		return
	}

	isTimeOut := false
	select {
	case <-time.NewTimer(12 * time.Millisecond).C:
	case result := <-timeout:
		if strings.Contains(result.Error(), "timeout") {
			isTimeOut = true
		}
	}
	if isTimeOut {
		t.Error("expected no timeout, but it did't")
	}
}
