package devicemanager

import (
	"log"

	"github.com/deepakkamesh/cm11"
)

type x10 struct {
	deviceNumber DeviceNumber
	cmd          chan DeviceData
	quit         chan struct{}
	err          chan error
	data         chan DeviceData
	tty          string
	cm11         *cm11.Device
	cm11Data     chan cm11.ObjState
}

func (m *x10) execute(data interface{}, address string) error {
	d, _ := data.(string)
	m.cm11.SendCommand(address[0:1], address[1:], d)
	log.Printf("cm11 executing %d on %s", d, address)
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
	m.cm11Data = make(chan cm11.ObjState)
	m.cm11 = cm11.New(m.tty, m.cm11Data)
	if err := m.cm11.Init(); err != nil {
		return err
	}
	go m.run()
	return nil
}

func (m *x10) Execute(action interface{}, object string) {
	m.cmd <- DeviceData{
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
		case cmd := <-m.cmd:
			if err := m.execute(cmd.Data, cmd.Object); err != nil {
				m.err <- err
			}
		case data := <-m.cm11Data:
			obj := data.HouseCode + data.DeviceCode // eg. C4
			m.data <- DeviceData{
				DeviceNumber: m.deviceNumber,
				Data:         data.FunctionCode,
				Object:       obj,
			}
		case <-m.quit:
			return
		}
	}
}
