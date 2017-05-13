/* Package devicemanager provides an abstraction layer over the devices
eg. X10 - on C4
texttospeech
http - get temp/weather
pandora - play electronica
*/
package devicemanager

import (
	"reflect"
	"regexp"

	"github.com/deepakkamesh/viki/devicemanager/device"
	"github.com/deepakkamesh/viki/objectmanager"
	"github.com/golang/glog"
)

// Define all the device types.
type DeviceId string

// DeviceData is used to send/recieve data with the device.
type DeviceData struct {
	DeviceId DeviceId    // Device number.
	Address  string      // Address of object.
	Data     interface{} // command to be executed or return value.
}

type DeviceSettings struct {
	Devices map[DeviceId]device.Device   // map of all the configured devices.
	Objects *objectmanager.ObjectManager // Objects
	Data    chan DeviceData              // channel to receive data from devices
	Err     chan error                   // channel to receive errors from devices
}

// New initializes a new device manager backend.
func New(o *objectmanager.ObjectManager) *DeviceSettings {
	glog.Infof("initializing device manager...")
	errChan := make(chan error, 10)       // Shared error channel.
	dataChan := make(chan DeviceData, 10) // Shared data channel.

	deviceSettings := &DeviceSettings{
		Devices: make(map[DeviceId]device.Device),
		Objects: o,
		Data:    dataChan,
		Err:     errChan,
	}

	// Call the initialization function for  devices.
	typ := reflect.TypeOf(&DeviceSettings{})
	for i := 0; i < typ.NumMethod(); i++ {
		if regexp.MustCompile("^NewDevice(.+)").MatchString(typ.Method(i).Name) {
			ret := reflect.ValueOf(&DeviceSettings{}).Method(i).Call(
				[]reflect.Value{
					reflect.ValueOf(dataChan),
					reflect.ValueOf(errChan),
					reflect.ValueOf(o),
				})
			devId := ret[0].Interface().(DeviceId)
			dev := ret[1].Interface().(device.Device)
			deviceSettings.Devices[devId] = dev
		}
	}

	return deviceSettings
}

// StartDeviceManager starts the devices .
func (m *DeviceSettings) StartDeviceManager() {

	// Start all the configured devices.
	for _, dev := range m.Devices {
		if err := dev.Start(); err != nil {
			glog.Errorf("Error starting device.  Failed with %s", err)
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
