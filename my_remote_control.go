package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

func (m *Viki) MyRemoteCode(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine userCode...")

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			name, _ := m.ObjectManager.GetObjectByAddress(got.Object)

			switch name {

			case "remote_1":
				if m.getMochadState("remote_1") == "On" {
					m.execObject("living light", "Off")
					m.execObject("dining light", "Off")
					m.execObject("patio light", "Off")
					m.execObject("tv light", "Off")
					break
				}
				m.execObject("living light", "On")
				m.execObject("dining light", "On")
			case "remote_2":
				m.execObject("bedroom light", m.getMochadState("remote_2"))
			case "remote_3":
				m.execObject("heater", m.getMochadState("remote_3"))
			}

		}
	}
}
