package devicemanager

import "log"

type deviceNumber int32

// Define all the device types.
const (
	Device_MISTERHOUSE deviceNumber = 1
	Device_PANDORA     deviceNumber = 2
)

// Interface implemented by all device types.
type device interface {
	Execute(string, string)
	On()
	Off()
	Start()
	Shutdown()
	GetErrorChan() <-chan error
}

type deviceCommand struct {
	deviceNumber deviceNumber
	action       string
	object       string
}

type deviceSettings struct {
	Devices map[deviceNumber]device
	quit    chan struct{}
	cmd     chan deviceCommand
}

// New initializes a new device manager backend.
func New() *deviceSettings {
	log.Printf("initializing device manager...")
	return &deviceSettings{
		Devices: map[deviceNumber]device{
			Device_MISTERHOUSE: &mh{
				ipPort: "1234",
				host:   "mh.utopia.com",
				cmd:    make(chan deviceCommand),
				quit:   make(chan struct{}),
				err:    make(chan error),
			},
			Device_PANDORA: &pandora{
				ip:   "1234",
				cmd:  make(chan deviceCommand),
				quit: make(chan struct{}),
				err:  make(chan error),
			},
		},
		quit: make(chan struct{}),
		cmd:  make(chan deviceCommand, 10),
	}
}

// StartDeviceManager starts the device manager backend.
func (m *deviceSettings) StartDeviceManager() {
	log.Printf("starting device manager...")

	// Start all the configured devices.
	m.Devices[Device_MISTERHOUSE].Start()
	m.Devices[Device_PANDORA].Start()
	go m.runDeviceManager()
}

// runDeviceManager runs the device manager loop.
func (m *deviceSettings) runDeviceManager() {
	for {
		select {
		case cmd := <-m.cmd:
			go m.Devices[cmd.deviceNumber].Execute(cmd.action, cmd.object)
		case <-m.quit:
			return
		case err := <-m.Devices[Device_PANDORA].GetErrorChan():
			log.Printf("Error executing Pandora %s", err)
		case err := <-m.Devices[Device_MISTERHOUSE].GetErrorChan():
			log.Printf("Error executing Misterhouse %s", err)
		}
	}
}

// ExecDeviceCommand executes the command on the specified device.
func (m *deviceSettings) ExecDeviceCommand(device deviceNumber, action string, object string) {
	m.cmd <- deviceCommand{
		deviceNumber: device,
		action:       action,
		object:       object,
	}
}

// shutdownDeviceManager shutdowns the device manager loop.
func (m *deviceSettings) ShutdownDeviceManager() {
	m.Devices[Device_MISTERHOUSE].Shutdown()
	m.Devices[Device_PANDORA].Shutdown()
	m.quit <- struct{}{}
}
