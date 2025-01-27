package tinygocapacitive

import (
	"machine"
	"time"
)

const (
	DEFAULT_SLEEP_TIME      time.Duration = 100000 * time.Nanosecond // 0.1 millisecond
	DEFAULT_TOUCH_THRESHOLD uint16        = 50000
	DEFAULT_MAX_ITERATIONS  uint          = 30
)

type Pad struct {
	adc machine.ADC

	// Time to sleep between checks while buffering.
	// Higher values may add notable delay to your main loop.
	// Lower values may not give the register time to change much.
	// defaults to DEFAULT_PAD_SLEEP_TIME.
	SleepTime time.Duration
	// Lower bound (between 0 and UINT16_MAX) to count a touch.
	// If this is too low, the floating/background noise will cause false HIGHs.
	TouchThreshold uint16
	// Number of times to check the value of the pin's register.
	// Sleeps for SleepTime prior to each check.
	// Returns early if a high enough value is found.
	MaxIterations uint
}

// Return a new Pad on the given pin, using default values.
// The given pin is configured as an ADC input and does not need to be initialized in advance.
func New(pin machine.Pin) *Pad {
	pin.Configure(machine.PinConfig{Mode: machine.PinInput})

	adc := machine.ADC{
		Pin: pin,
	}

	adc.Configure(machine.ADCConfig{})

	return &Pad{
		adc:            adc,
		SleepTime:      DEFAULT_SLEEP_TIME,
		TouchThreshold: DEFAULT_TOUCH_THRESHOLD,
		MaxIterations:  DEFAULT_MAX_ITERATIONS,
	}
}

// Returns whether or not the pad is currently touched (HIGH), smoothing over fluctuations and floating values.
func (p *Pad) Get() bool {
	for i := 0; i < 30; i++ {
		time.Sleep(p.SleepTime)
		if p.adc.Get() > p.TouchThreshold {
			return true
		}
	}
	return false
}
