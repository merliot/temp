package temp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/merliot/dean"
	"github.com/merliot/device"
	"github.com/merliot/device/dht"
)

var targets = []string{"demo", "nano-rp2040", "wioterminal"}

type Temp struct {
	*device.Device
	Dht     dht.Dht
	prevDht dht.Dht
}

type MsgUpdate struct {
	Path        string
	Temperature float32
	Humidity    float32
}

func New(id, model, name string) dean.Thinger {
	fmt.Println("NEW TEMP")
	return &Temp{
		Device: device.New(id, model, name, fs, targets).(*device.Device),
	}
}

func (t *Temp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.API(w, r, t)
}

func (t *Temp) save(msg *dean.Msg) {
	msg.Unmarshal(t).Broadcast()
}

func (t *Temp) getState(msg *dean.Msg) {
	t.Path = "state"
	msg.Marshal(t).Reply()
}

func (t *Temp) update(msg *dean.Msg) {
	msg.Unmarshal(&t.Dht).Broadcast()
}

func (t *Temp) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"state":     t.save,
		"get/state": t.getState,
		"update":    t.update,
	}
}

func (t *Temp) parseParams() {
	t.Dht.Sensor = t.ParamFirstValue("sensor")
	t.Dht.Gpio = t.ParamFirstValue("gpio")
	t.Dht.Configure()
}

func (t *Temp) Setup() {
	t.Device.Setup()
	t.parseParams()
}

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
