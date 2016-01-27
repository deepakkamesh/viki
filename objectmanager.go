package viki

import "viki/devicemanager"

/* Object manager reads the configuration and maps the objects with the underlying
device manager. It also maintains state of each object.
*/

type Object struct {
	Name    string                        // Human readable name.
	Address string                        // Address of device. Optional.
	State   interface{}                   // State of object.
	DevNo   devicemanager.DeviceNumber    // Device number.
	DevMgr  *devicemanager.DeviceSettings // DeviceMaanger.
}

func InitObject(name string, address string, deviceNo devicemanager.DeviceNumber, deviceMgr *devicemanager.DeviceSettings) *Object {
	return &Object{
		Name:    name,
		Address: address,
		DevNo:   deviceNo,
		DevMgr:  deviceMgr,
	}
}

func (m *Object) Execute(data interface{}) {
	m.DevMgr.Devices[m.DevNo].Execute(data, m.Address)
	m.State = data
}
