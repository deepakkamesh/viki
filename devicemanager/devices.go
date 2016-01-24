/* Package devicemanager provides an abstraction layer over the devices
eg. X10 - on C4
texttospeech
http - get temp/weather
pandora - play electronica
*/
package devicemanager

import "log"

type DeviceNumber int32

// Define all the device types.
const (
	Device_MISTERHOUSE DeviceNumber = 1
	Device_HTTPHANDLER DeviceNumber = 2
)

// Interface implemented by all device types.
type device interface {
	Execute(string, string)
	On()
	Off()
	Start()
	Shutdown()
}

// Command structure for devices.
type DeviceCommand struct {
	DeviceNumber DeviceNumber
	Action       string
	Object       string
}

// Return Data from devices.
type DeviceData struct {
	DeviceNumber DeviceNumber
	Data         interface{}
	Object       string
}

type DeviceSettings struct {
	Devices map[DeviceNumber]device
	quit    chan struct{}
	cmd     chan DeviceCommand
	Data    chan DeviceData
	Err     chan error
}

// New initializes a new device manager backend.
func New() *DeviceSettings {
	log.Printf("initializing device manager...")
	errChan := make(chan error, 10)       // Shared error channel.
	dataChan := make(chan DeviceData, 10) // Shared data channel.

	return &DeviceSettings{
		Devices: map[DeviceNumber]device{
			Device_MISTERHOUSE: &mh{
				deviceNumber: Device_MISTERHOUSE,
				cmd:          make(chan DeviceCommand),
				quit:         make(chan struct{}),
				err:          errChan,
				data:         dataChan,
			},
			Device_HTTPHANDLER: &HttpHandler{
				deviceNumber: Device_HTTPHANDLER,
				cmd:          make(chan DeviceCommand),
				quit:         make(chan struct{}),
				err:          errChan,
				data:         dataChan,
			},
		},
		quit: make(chan struct{}),
		cmd:  make(chan DeviceCommand, 10),
		Data: dataChan,
		Err:  errChan,
	}
}

// runDeviceManager runs the device manager loop.
func (m *DeviceSettings) runDeviceManager() {
	for {
		select {
		case cmd := <-m.cmd:
			go m.Devices[cmd.DeviceNumber].Execute(cmd.Action, cmd.Object)
		case <-m.quit:
			return
		}
	}
}

// StartDeviceManager starts the device manager backend.
func (m *DeviceSettings) StartDeviceManager() {

	// Start all the configured devices.
	for _, dev := range m.Devices {
		dev.Start()
	}
	go m.runDeviceManager()
}

// ExecDeviceCommand executes the command on the specified device.
func (m *DeviceSettings) ExecDeviceCommand(device DeviceNumber, action string, object string) {
	m.cmd <- DeviceCommand{
		DeviceNumber: device,
		Action:       action,
		Object:       object,
	}
}

// shutdownDeviceManager shutdowns the device manager loop.
func (m *DeviceSettings) ShutdownDeviceManager() {

	// Stop all the configured devices.
	for _, dev := range m.Devices {
		dev.Shutdown()
	}
	//Shutdown device manager.
	m.quit <- struct{}{}
}
