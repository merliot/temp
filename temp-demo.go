//go:build !tinygo && !rpi

package temp

import (
	"math"
	"math/rand"
	"time"

	"github.com/merliot/dean"
)

func randomValue(mean, stddev float64) float32 {
	u1 := rand.Float64()
	u2 := rand.Float64()

	// Box-Muller transform to generate normally distributed random numbers
	z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2)
	value := mean + stddev*z0

	// round to 1 dec place
	return float32(math.Round(value*10) / 10)
}

func (t *Temp) run(i *dean.Injector) {
	var msg dean.Msg
	var d = &t.Dht
	var prev = &t.prevDht

	d.Temperature = randomValue(75.0, 0.05)
	d.Humidity = randomValue(36.4, 0.05)

	if d.Temperature != prev.Temperature ||
		d.Humidity != prev.Humidity {
		var update = MsgUpdate{
			Path:        "update",
			Temperature: d.Temperature,
			Humidity:    d.Humidity,
		}
		i.Inject(msg.Marshal(update))
		prev.Temperature = d.Temperature
		prev.Humidity = d.Humidity
	}
}

func (t *Temp) Run(i *dean.Injector) {
	rand.Seed(time.Now().UnixNano())
	for {
		t.run(i)
		// limit reads to every 2 seconds
		time.Sleep(2 * time.Second)
	}
}
