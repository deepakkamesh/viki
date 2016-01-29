/* Example stub code for a device driver.
Steps to add a new device driver.
1. Make a copy of this file and update.
2. Add the device to devices.go
*/
package devicemanager

import "log"

// Unique Device Id. Usually  same as device name.
const Device_ANYDEVICE DeviceId = "anydevice"

type anydev struct {
	deviceId DeviceId
	in       chan DeviceData
	quit     chan struct{}
	err      chan error
	out      chan DeviceData
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

// DONOTCHANGE.
// Start initiates the device.
func (m *anydev) Start() error {
	log.Printf("starting device [name]...")
	// Set any required parameters using flag.
	go m.run()
	return nil
}

// DONOTCHANGE.
// Execute queues up the requested command to the channel.
func (m *anydev) Execute(action interface{}, object string) {
	m.in <- DeviceData{
		Data:   action,
		Object: object,
	}
}

// DONOTCHANGE.
// Shutdown terminates the device processing.
func (m *anydev) Shutdown() {
	m.quit <- struct{}{}
}

// DONOTCHANGE.
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
