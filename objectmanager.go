package viki

import "viki/devicemanager"

/* Object manager reads the configuration and maps the objects with the underlying
device manager. It also maintains state of each object.
*/

type Object struct {
	Address string               // Address of device. Optional.
	device  devicemanager.Device // underlying device driver.
	State   interface{}          // State of object.
}

func InitObject(address string, device devicemanager.Device) *Object {
	return &Object{
		Address: address,
		device:  device,
	}
}

func (m *Object) Execute(data interface{}) {
	m.device.Execute(data, m.Address)
	m.State = data
}
