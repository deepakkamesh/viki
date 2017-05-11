/* Example stub code for a device driver.
Steps to add a new device driver.
1. Make a copy of this file and update.
2. Add the device to devices.go
*/
package devicemanager

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager/device"
	"github.com/deepakkamesh/viki/objectmanager"
)

// Unique Device Id. Usually  same as device name.
const Device_ANYDEVICE DeviceId = "anydevice"

type anydev struct {
	in   chan DeviceData
	quit chan struct{}
	err  chan error
	out  chan DeviceData
	om   *objectmanager.ObjectManager
}

// NewDevice<deviceName> returns a new and initialized anydevice.
// The function needs to start with NewDevice* for device manager
// to recognize this as a initializing function. Anything else
// is ignored. It also has to implement the Device Interface.
func (m *DeviceSettings) newDeviceAnyDevice(out chan DeviceData, err chan error, om *objectmanager.ObjectManager) (DeviceId, device.Device) {
	return Device_ANYDEVICE, &anydev{
		in:   make(chan DeviceData, 10), // Input channel, typically buffered.
		quit: make(chan struct{}),       // Quit.
		err:  err,                       // Common error channel.
		out:  out,                       // Channel to send out data.
		om:   om,
	}
}

func (m *anydev) execute(data interface{}, object string) error {

	// Assert the command data depending on device.
	d, _ := data.(string)
	log.Printf("anydevice: executing %d on %s", d, object)
	return nil
}

// On turns off the device.
func (m *anydev) On() {
	log.Printf("Turn off")
}

// Off turns off the device.
func (m *anydev) Off() {
	log.Printf("Turn off")
}

// Start initiates the device.
func (m *anydev) Start() error {
	log.Printf("starting device [name]...")
	// Set any required parameters using flag.
	go m.run()
	return nil
}

// Execute queues up the requested command to the channel.
func (m *anydev) Execute(action interface{}, object string) {
	m.in <- DeviceData{
		Data:   action,
		Object: object,
	}
}

// Shutdown terminates the device processing.
func (m *anydev) Shutdown() {
	m.quit <- struct{}{}
}

// run is the main processing loop for the device driver.
func (m *anydev) run() {
	for {
		select {
		case in := <-m.in:
			if err := m.execute(in.Data, in.Object); err != nil {
				m.err <- err
			}
		case <-m.quit:
			return
		}
	}
}
