# Tiny Go Capacitive Touch Pads

A smoothing wrapper for the Adafruit Circuit Playground Express's capacitive touchpads. If a touch is noticed anywhere within a fraction of a millisecond from the Get() call, any flickering or floating will be papered over. This allows you to actually use the capacitive pads as just booleans.


## Example Usage

```go
package main

import (
	"machine"
	"time"
)

var p *Pad = New(machine.A1)

func init() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func main() {

	for {
		machine.LED.Set(p.Get())

		time.Sleep(100 * time.Millisecond)
	}
}
```


## Why?

Because I could not get anything else to work for the CPX.

When reading the capacitive touchpads (A1-A7), the results flip-flop (presumably meaning the pins are "floating"). Treating the pins as digital would cause them to "blink" true every so often. Treating the pins as analog showed why: they hover around ~20-30k. Easy fix, pull them down, right? Or pull them up and invert the result? Yeah, no dice. Pulled up, down, or left floating, they would alternate between UINT16_MAX and ~0 when contacted, causing the result to flicker when touched and then float when not touched.

- TinyGo's touch/capacitive library did not work; it always returned the same number for each pad, no matter if the pad was being touched (ex: A1 was something like 354).

- TinyGo's MakeyButton didn't work because of the aforementioned 

### So How Does this Library Fix That?

Basically, it buffers for the existence of a high value (touch) for a fraction of a millisecond on read, thereby alleviating both the floating of the pins and the flickering between HIGH and LOW that they detect. If any values over a reasonably high threshold are found (which they will, as touching causes it to flip between UINT16_MAX and 0), count the whole Get() as HIGH.