package devicemanager

import (
	"flag"
	"log"

	"github.com/deepakkamesh/cm11"
)

// Unique Device Number.
const Device_X10 DeviceId = "x10"

type x10 struct {
	deviceId DeviceId
	in       chan DeviceData
	quit     chan struct{}
	err      chan error
	out      chan DeviceData
	cm11     *cm11.Device
	cm11Data chan cm11.ObjState
}

func (m *x10) execute(data interface{}, address string) error {
	d, _ := data.(string)
	m.cm11.SendCommand(address[0:1], address[1:], d)
	log.Printf("cm11 executing %s on %s", d, address)
	return nil
}

func (m *x10) On() {
	log.Printf("Turn on x10")
}
func (m *x10) Off() {
	log.Printf("Turn off x10")
}

func (m *x10) Start() error {
	log.Printf("starting device cm11...")
	fl := flag.Lookup("x10_tty")
	tty := fl.Value.String()
	m.cm11Data = make(chan cm11.ObjState)
	m.cm11 = cm11.New(tty, m.cm11Data)
	if err := m.cm11.Init(); err != nil {
		return err
	}
	go m.run()
	return nil
}

func (m *x10) Execute(action interface{}, object string) {
	m.in <- DeviceData{
		Data:   action,
		Object: object,
	}
}
func (m *x10) Shutdown() {
	m.quit <- struct{}{}
}

func (m *x10) run() {
	for {
		select {
		case in := <-m.in:
			if err := m.execute(in.Data, in.Object); err != nil {
				m.err <- err
			}
		case data := <-m.cm11Data:
			obj := data.HouseCode + data.DeviceCode // eg. C4
			m.out <- DeviceData{
				DeviceId: m.deviceId,
				Data:     data.FunctionCode,
				Object:   obj,
			}
		case <-m.quit:
			return
		}
	}
}
