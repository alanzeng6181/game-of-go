package gamelogic

import (
	"errors"
	"sync/atomic"
	"time"
)

type ClockConfig struct {
	base       time.Duration
	ticking    time.Duration
	violations uint8
}

type Clock struct {
	timeout chan struct {
		Stone
		error
	}
	//e.g. 30min base, 30secons ticker, 3 violations would mean that when 30min time
	// is used up, there will be a countdown of 30 seconds for each move. A player can
	// hit maximum of 3 violations, after which the player automatically loses.
	ClockConfig
	input   <-chan Stone
	started atomic.Value
}

func NewClock(clockConfig ClockConfig, input <-chan Stone) (*Clock, error) {
	if clockConfig.base == time.Duration(0) && clockConfig.ticking == time.Duration(0) {
		return nil, errors.New("at least the base or ticking has be non zero")
	}
	var started atomic.Value
	started.Store(false)
	return &Clock{
		timeout: make(chan struct {
			Stone
			error
		}, 1),
		ClockConfig: clockConfig,
		input:       input,
		started:     started,
	}, nil
}

func (clock *Clock) GetTimeout() (<-chan struct {
	Stone
	error
}, error) {
	if !clock.started.Load().(bool) {
		return nil, errors.New("clock is not started")
	}

	return (<-chan struct {
		Stone
		error
	})(clock.timeout), nil
}

func (clock *Clock) Start() {
	if !clock.started.CompareAndSwap(false, true) {
		return
	}
	go func() {
		var timer1 *PausableTimer
		var timer2 *PausableTimer
		var countDown1 uint8 = 0
		var countDown2 uint8 = 0
		expectedColor := Black
		if clock.base > time.Duration(0) {
			timer1 = NewPausableTimer(clock.base)
			countDown1 = 0
			timer2 = NewPausableTimer(clock.base)
			countDown2 = 0
		} else {
			timer1 = NewPausableTimer(clock.ticking)
			countDown1 = 1
			timer2 = NewPausableTimer(clock.ticking)
			countDown2 = 1
		}
		defer func() {
			timer1.Stop()
			timer2.Stop()
		}()
		timer1C := timer1.C()
		timer2C := timer2.C()
		for {
			select {
			case stone := <-clock.input:
				if expectedColor != stone {
					clock.timeout <- struct {
						Stone
						error
					}{stone, errors.New("expected different color")}
					return
				} else if stone == Black {
					if countDown1 == 0 {
						timer1.Pause()
					} else {
						timer1.Reset(clock.ticking)

					}
				} else {
					if countDown2 == 0 {
						timer2.Pause()
					} else {
						timer2.Reset(clock.ticking)
					}
				}
				expectedColor = expectedColor.OppositeColor()
			case <-timer1C:
				if countDown1 > clock.violations {
					clock.timeout <- struct {
						Stone
						error
					}{Black, errors.New("timeout")}
					close(clock.timeout)
					return
				} else {
					countDown1++
					timer1.Reset(clock.ticking)
				}
			case <-timer2C:
				if countDown2 > clock.violations {
					clock.timeout <- struct {
						Stone
						error
					}{Black, errors.New("timeout")}
					return
				} else {
					countDown2++
					timer2.Reset(clock.ticking)
				}
			}
		}
	}()
}

type PausableTimer struct {
	timer     *time.Timer
	lastStart *time.Time
	lastPause *time.Time
	duration  time.Duration
}

func NewPausableTimer(duration time.Duration) *PausableTimer {
	now := time.Now()
	return &PausableTimer{timer: time.NewTimer(duration), duration: duration, lastStart: &now}
}

func (pt *PausableTimer) Stop() {
	pt.timer.Stop()
}

func (pt *PausableTimer) Pause() {
	now := time.Now()
	pt.duration = pt.duration - now.Sub(*pt.lastStart)
	pt.lastPause = &now
}

func (pt *PausableTimer) Resume() {
	now := time.Now()
	pt.duration = pt.duration - now.Sub(*pt.lastStart)
	*pt.lastStart = now
	pt.timer.Reset(pt.duration)
}

func (pt *PausableTimer) Reset(duration time.Duration) {
	pt.timer.Reset(duration)
	pt.duration = duration
	*pt.lastStart = time.Now()
	pt.lastPause = nil
}

func (pt *PausableTimer) C() <-chan time.Time {
	return pt.timer.C
}
