package objectmanager

import (
	"strings"

	"github.com/deepakkamesh/viki/devicemanager/device"
)

//ObjectManager manages the objects.
type ObjectManager struct {
	Objects []*Object
}

func New() *ObjectManager {
	return &ObjectManager{}
}

func (m *ObjectManager) AddObject(name string, address string, device device.Device, tags []string) error {
	m.Objects = append(m.Objects, NewObject(name, address, device, tags))
	return nil
}

func (m *ObjectManager) GetObjectByName(name string) (string, *Object) {
	for _, o := range m.Objects {
		if o.Name == name {
			return o.Address, o
		}
	}
	return "", nil
}

func (m *ObjectManager) GetObjectByAddress(address string) (string, *Object) {
	for _, o := range m.Objects {
		if o.Address == address {
			return o.Name, o
		}
	}
	return "", nil
}

/* Object manager reads the configuration and maps the objects with the underlying
device manager. It also maintains state of each object.
*/
type Object struct {
	Address string        // Address of device. Optional.
	Name    string        // Human readable name of device.
	device  device.Device // underlying device driver.
	State   interface{}   // State of object.
	Tags    []string      // Tags associated with object.
}

func NewObject(name string, address string, device device.Device, tags []string) *Object {
	return &Object{
		Name:    name,
		Address: address,
		device:  device,
		Tags:    tags,
	}
}

// Execute calls the underlying device driver to execute command.
func (m *Object) Execute(data interface{}) error {
	m.SetState(data)
	// If object is associated with a device send an exec command to it.
	if m.device != nil {
		m.device.Execute(data, m.Address)
	}
	return nil
}

// SetState changes state of object.
func (m *Object) SetState(data interface{}) {
	m.State = data
}

//getTag returns the tag that matches the object substr.
func (m *Object) GetTag(tag string) (string, bool) {
	for _, k := range m.Tags {
		if strings.Contains(k, tag) {
			return k, true
		}
	}
	return "", false
}

// CheckTag returns true if the tag exists on object.
func (m *Object) CheckTag(tag string) bool {
	for _, k := range m.Tags {
		if k == tag {
			return true
		}
	}
	return false
}
