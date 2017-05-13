package viki

import (
	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/golang/glog"
)

func (m *Viki) MyRemoteCode(c chan devicemanager.DeviceData) {

	glog.Info("starting user routine MyRemoteCode...")
	defer glog.Info("shutting down user routine MyRemoteCode")

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			name, _ := m.ObjectManager.GetObjectByAddress(got.Address)

			switch name {
			case "remote_1":
				if m.getMochadState("remote_1") == "On" {
					m.Do("living light", "Off")
					m.Do("dining light", "Off")
					m.Do("patio light", "Off")
					m.Do("tv light", "Off")
					continue
				}
				m.Do("living light", "On")
				m.Do("dining light", "On")
			case "remote_2":
				m.Do("bedroom light", m.getMochadState("remote_2"))
			case "remote_3":
				m.Do("heater", m.getMochadState("remote_3"))
			}

		}
	}
}
