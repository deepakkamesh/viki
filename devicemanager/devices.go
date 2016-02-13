/* Package devicemanager provides an abstraction layer over the devices
eg. X10 - on C4
texttospeech
http - get temp/weather
pandora - play electronica
*/
package devicemanager

import "log"

// Define all the device types.
type DeviceId string

// Interface implemented by all device types.
type Device interface {
	Execute(interface{}, string)
	On()
	Off()
	Start() error
	Shutdown()
}

// DeviceData is used to communicate with the device.
type DeviceData struct {
	DeviceId DeviceId    // Device number.
	Object   string      // Address of object.
	Data     interface{} // command to be executed or return value.
}

type DeviceSettings struct {
	Devices map[DeviceId]Device
	Data    chan DeviceData
	Err     chan error
}

// New initializes a new device manager backend.
func New() *DeviceSettings {
	log.Printf("initializing device manager...")
	errChan := make(chan error, 10)       // Shared error channel.
	dataChan := make(chan DeviceData, 10) // Shared data channel.

	// Add new devices here.
	return &DeviceSettings{
		Devices: map[DeviceId]Device{
			Device_X10: &x10{
				deviceId: Device_X10,
				in:       make(chan DeviceData, 10),
				quit:     make(chan struct{}),
				err:      errChan,
				out:      dataChan,
			},
			Device_HTTPHANDLER: &httphandler{
				deviceId: Device_HTTPHANDLER,
				in:       make(chan DeviceData, 10),
				quit:     make(chan struct{}),
				err:      errChan,
				out:      dataChan,
			},
			Device_SPEAKER: &speaker{
				deviceId: Device_SPEAKER,
				in:       make(chan DeviceData, 10),
				quit:     make(chan struct{}),
				err:      errChan,
				out:      dataChan,
			},
			Device_MOCHAD: &mochad{
				deviceId: Device_MOCHAD,
				in:       make(chan DeviceData, 10),
				quit:     make(chan struct{}),
				err:      errChan,
				out:      dataChan,
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
