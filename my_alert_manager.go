// +build ignore

package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

func (m *Viki) myAlertManager(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine Alert Manager...")

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			name, obj := m.getObject(got.Object)

			// Alerts when we are at home.
			if m.getModeState("mode away") == "Off" {
				st := m.getMochadState(name)
				// If door is opened.
				if st == "Open" && obj.checkTag("door") {
					m.execObject("speaker", "Warning "+name+" is open")
					continue
				}
				// if motion sensor backyard and door is not open.
				if st == "On" && name == "backyard_ms1" && m.getMochadState("backyard door") != "Open" {
					m.execObject("speaker", "Warning backyard motion sensor activated ")
					continue
				}
				// if motion sensor garage and door is not open.
				if st == "On" && name == "garage_ms1" && m.getMochadState("garage door") != "Open" {
					m.execObject("speaker", "Warning garage motion sensor activated ")
					continue
				}
			}

		}

	}
}
