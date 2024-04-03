//go:build tinygo

package temp

import (
	"embed"
	"time"

	"github.com/merliot/dean"
)

var fs embed.FS

func (t *Temp) run(i *dean.Injector) {
	var msg dean.Msg
	var d = &t.Dht
	var prev = &t.prevDht

	err := d.Read()
	if err != nil {
		println("read err", err.Error())
		return
	}

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
	for {
		t.run(i)
		// limit reads to every 2 seconds
		time.Sleep(2 * time.Second)
	}
}
