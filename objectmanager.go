package viki

import "github.com/deepakkamesh/viki/devicemanager"

/* Object manager reads the configuration and maps the objects with the underlying
device manager. It also maintains state of each object.
*/

type Object struct {
	Address string               // Address of device. Optional.
	device  devicemanager.Device // underlying device driver.
	State   interface{}          // State of object.
	Tags    []string             // Tags associated with object.
}

func InitObject(address string, device devicemanager.Device, tags []string) *Object {
	return &Object{
		Address: address,
		device:  device,
		Tags:    tags,
	}
}

// Execute calls the underlying device driver to execute command.
func (m *Object) execute(data interface{}) {
	if m.device != nil {
		m.device.Execute(data, m.Address)
	}
	m.setState(data)
}

// SetState changes state of object.
func (m *Object) setState(data interface{}) {
	m.State = data
}

// CheckTag returns true if the tag exists on object.
func (m *Object) checkTag(tag string) bool {
	for _, k := range m.Tags {
		if k == tag {
			return true
		}
	}
	return false
}
