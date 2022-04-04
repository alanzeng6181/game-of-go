package gamelogic

import (
	"errors"
	"sync/atomic"
	"time"
)

type ClockConfig struct {
	TimeLimit     time.Duration
	OverTime      time.Duration
	OverTimeCount uint8
}

func MakeClockConfig(base time.Duration, ticking time.Duration, violations uint8) ClockConfig {
	return ClockConfig{base, ticking, violations}
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
	input      <-chan Stone
	started    atomic.Value
	timer1     *PausableTimer
	timer2     *PausableTimer
	countDown1 uint8
	countDown2 uint8
}

func NewClock(clockConfig ClockConfig, input <-chan Stone) (*Clock, error) {
	if clockConfig.TimeLimit == time.Duration(0) && clockConfig.OverTime == time.Duration(0) {
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
	return (<-chan struct {
		Stone
		error
	})(clock.timeout), nil
}

func (clock *Clock) Start() {
	if !clock.started.CompareAndSwap(false, true) {
		return
	}
	if clock.TimeLimit > time.Duration(0) {
		clock.timer1 = NewPausableTimer(clock.TimeLimit)
		clock.countDown1 = 0
		clock.timer2 = NewPausableTimer(clock.TimeLimit)
		clock.countDown2 = 0
	} else {
		clock.timer1 = NewPausableTimer(clock.OverTime)
		clock.countDown1 = 1
		clock.timer2 = NewPausableTimer(clock.OverTime)
		clock.countDown2 = 1
	}
	go func() {
		expectedColor := Black
		defer func() {
			clock.timer1.Stop()
			clock.timer2.Stop()
		}()
		timer1C := clock.timer1.C()
		timer2C := clock.timer2.C()
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
					if clock.countDown1 == 0 {
						clock.timer1.Pause()
					} else {
						clock.timer1.Reset(clock.OverTime)

					}
				} else {
					if clock.countDown2 == 0 {
						clock.timer2.Pause()
					} else {
						clock.timer2.Reset(clock.OverTime)
					}
				}
				expectedColor = expectedColor.OppositeColor()
			case <-timer1C:
				if clock.countDown1 > clock.OverTimeCount {
					clock.timeout <- struct {
						Stone
						error
					}{Black, errors.New("timeout")}
					close(clock.timeout)
					return
				} else {
					clock.countDown1++
					clock.timer1.Reset(clock.OverTime)
				}
			case <-timer2C:
				if clock.countDown2 > clock.OverTimeCount {
					clock.timeout <- struct {
						Stone
						error
					}{Black, errors.New("timeout")}
					return
				} else {
					clock.countDown2++
					clock.timer2.Reset(clock.OverTime)
				}
			}
		}
	}()
}

func (clock *Clock) ToClockState() ClockState {
	now := time.Now()

	return ClockState{
		Timestamp:      now,
		BlackOTCount:   clock.countDown1,
		WhiteOTCount:   clock.countDown2,
		BlackTimerRead: uint32(clock.timer1.TimerRead().Seconds()),
		WhiteTimerRead: uint32(clock.timer2.TimerRead().Seconds()),
		TimeLimit:      uint32(clock.TimeLimit.Seconds()),
		OTCout:         clock.OverTimeCount,
	}
}

type PausableTimer struct {
	originalDuration time.Duration
	originalStart    time.Time
	timer            *time.Timer
	lastStart        time.Time
	lastPause        *time.Time
	duration         time.Duration
}

func NewPausableTimer(duration time.Duration) *PausableTimer {
	now := time.Now()
	return &PausableTimer{timer: time.NewTimer(duration), duration: duration, originalDuration: duration, originalStart: now, lastStart: now}
}

func (pt *PausableTimer) Stop() {
	pt.timer.Stop()
}

func (pt *PausableTimer) Pause() {
	now := time.Now()
	pt.duration = pt.duration - now.Sub(pt.lastStart)
	pt.lastPause = &now
}

func (pt *PausableTimer) Resume() {
	now := time.Now()
	pt.duration = pt.duration - now.Sub(pt.lastStart)
	pt.lastStart = now
	pt.timer.Reset(pt.duration)
}

func (pt *PausableTimer) Reset(duration time.Duration) {
	pt.originalDuration = duration
	now := time.Now()
	pt.originalStart = now
	pt.timer.Reset(duration)
	pt.duration = duration
	pt.lastStart = now
	pt.lastPause = nil
}

func (pt *PausableTimer) C() <-chan time.Time {
	return pt.timer.C
}

func (pt *PausableTimer) TimerRead() time.Duration {
	return time.Now().Sub(pt.originalStart)
}
