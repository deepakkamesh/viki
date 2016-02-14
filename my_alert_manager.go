package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

func (m *Viki) alertManager(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine Alert Manager...")

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			d, _ := got.Data.(string)
			if got.DeviceId == "mochad" {
				if obj := m.GetObject(got.Object); d == "Open" && obj != nil && obj.CheckTag("door") {
					door := m.GetObjectName(obj.Address)
					m.ExecObject("speaker", "Warning "+door+" is open")

				}

			}

		}
	}
}