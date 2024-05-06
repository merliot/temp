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

type Record [2]float32

type Temp struct {
	*device.Device
	Dht     dht.Dht
	History []Record
}

type MsgUpdate struct {
	Path        string
	Temperature float32
	Humidity    float32
}

func New(id, model, name string) dean.Thinger {
	fmt.Println("NEW TEMP\r")
	return &Temp{
		Device:  device.New(id, model, name, fs, targets).(*device.Device),
		History: []Record{},
	}
}

func (t *Temp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.API(w, r, t)
}

func (t *Temp) save(pkt *dean.Packet) {
	pkt.Unmarshal(t).Broadcast()
}

func (t *Temp) getState(pkt *dean.Packet) {
	t.Path = "state"
	pkt.Marshal(t).Reply()
}

func (t *Temp) addRecord() {
	if len(t.History) >= 60 {
		// Remove the oldest
		t.History = t.History[1:]
	}
	// Add the new
	r := Record{t.Dht.Temperature, t.Dht.Humidity}
	t.History = append(t.History, r)
}

func (t *Temp) update(pkt *dean.Packet) {
	pkt.Unmarshal(&t.Dht).Broadcast()
	t.addRecord()
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

func (t *Temp) minute(i *dean.Injector) {
	var pkt dean.Packet
	var d = &t.Dht

	err := d.Read()
	if err != nil {
		println("read err", err.Error())
		return
	}

	var update = MsgUpdate{
		Path:        "update",
		Temperature: d.Temperature,
		Humidity:    d.Humidity,
	}
	i.Inject(pkt.Marshal(update))
}

func (t *Temp) Run(i *dean.Injector) {
	for {
		t.minute(i)
		time.Sleep(60 * time.Second)
	}
}
