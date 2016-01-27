/* Package devicemanager provides an abstraction layer over the devices
eg. X10 - on C4
texttospeech
http - get temp/weather
pandora - play electronica
*/
package devicemanager

import "log"

// Predefined States.
type ObjState int32

const (
	State_ON     ObjState = 1
	State_OFF    ObjState = 2
	State_MOTION ObjState = 3
	State_STILL  ObjState = 4
)

// Define all the device types.
type DeviceNumber int32

const (
	Device_X10         DeviceNumber = 1
	Device_HTTPHANDLER DeviceNumber = 2
)

// Interface implemented by all device types.
type device interface {
	Execute(interface{}, string)
	On()
	Off()
	Start() error
	Shutdown()
}

type DeviceData struct {
	DeviceNumber DeviceNumber // Device number.
	Data         interface{}  // command to be executed or return value.
	Object       string       // Address of object.
}

type DeviceSettings struct {
	Devices map[DeviceNumber]device
	Data    chan DeviceData
	Err     chan error
}

// New initializes a new device manager backend.
func New() *DeviceSettings {
	// TODO: Read device config from file.
	log.Printf("initializing device manager...")
	errChan := make(chan error, 10)       // Shared error channel.
	dataChan := make(chan DeviceData, 10) // Shared data channel.

	// Add new devices here.
	return &DeviceSettings{
		Devices: map[DeviceNumber]device{
			Device_X10: &x10{
				deviceNumber: Device_X10,
				cmd:          make(chan DeviceData, 10),
				quit:         make(chan struct{}),
				err:          errChan,
				data:         dataChan,
				tty:          "/dev/ttyUSB0",
			},
			Device_HTTPHANDLER: &HttpHandler{
				deviceNumber: Device_HTTPHANDLER,
				cmd:          make(chan DeviceData, 10),
				quit:         make(chan struct{}),
				err:          errChan,
				data:         dataChan,
			},
		},
		Data: dataChan,
		Err:  errChan,
	}
}

// StartDeviceManager starts the devices .
func (m *DeviceSettings) StartDeviceManager() {

	// Start all the configured devices.
	for _, dev := range m.Devices {
		if err := dev.Start(); err != nil {
			log.Printf("Error starting device.  Failed with %s", err)
		}
	}
}

// shutdownDeviceManager shutdowns the device manager loop.
func (m *DeviceSettings) ShutdownDeviceManager() {

	// Stop all the configured devices.
	for _, device := range m.Devices {
		device.Shutdown()
	}
}
