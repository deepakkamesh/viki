package devicemanager

import (
	"flag"

	"github.com/deepakkamesh/cm11"
	"github.com/deepakkamesh/viki/devicemanager/device"
	"github.com/deepakkamesh/viki/objectmanager"
	"github.com/golang/glog"
)

const Device_X10 DeviceId = "x10"

type x10 struct {
	in       chan DeviceData
	quit     chan struct{}
	err      chan error
	out      chan DeviceData
	om       *objectmanager.ObjectManager
	cm11     *cm11.Device
	cm11Data chan cm11.ObjState
	cm11Err  chan error
}

// Function to initialize the device.
// Function called by devicemanager to initialize the device
func (m *DeviceSettings) NewDeviceX10(out chan DeviceData, err chan error, om *objectmanager.ObjectManager) (DeviceId, device.Device) {
	return Device_X10, &x10{
		in:       make(chan DeviceData, 10),
		quit:     make(chan struct{}),
		err:      err,
		out:      out,
		om:       om,
		cm11Data: make(chan cm11.ObjState),
		cm11Err:  make(chan error),
	}
}

func (m *x10) execute(data interface{}, address string) error {
	d, _ := data.(string)
	m.cm11.SendCommand(address[0:1], address[1:], d)
	glog.Infof("cm11 executing %s on %s", d, address)
	return nil
}

func (m *x10) On() {
	glog.Infof("Turn on x10")
}
func (m *x10) Off() {
	glog.Infof("Turn off x10")
}

func (m *x10) Start() error {
	glog.Infof("starting device cm11...")
	tty := flag.Lookup("x10_tty").Value.String()
	m.cm11 = cm11.New(tty, m.cm11Data, m.cm11Err)
	if err := m.cm11.Init(); err != nil {
		return err
	}
	go m.run()
	return nil
}

func (m *x10) Execute(action interface{}, object string) {
	m.in <- DeviceData{
		Data:    action,
		Address: object,
	}
}
func (m *x10) Shutdown() {
	m.quit <- struct{}{}
}

func (m *x10) run() {
	for {
		select {
		case in := <-m.in:
			if err := m.execute(in.Data, in.Address); err != nil {
				m.err <- err
			}
		case data := <-m.cm11Data:
			obj := data.HouseCode + data.DeviceCode // eg. C4
			m.out <- DeviceData{
				DeviceId: Device_X10,
				Data:     data.FunctionCode,
				Address:  obj,
			}
		case err := <-m.cm11Err:
			m.err <- err
		case <-m.quit:
			return
		}
	}
}
