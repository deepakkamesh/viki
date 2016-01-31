package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

func (m *Viki) modeSleep(in chan devicemanager.DeviceData) {

	log.Printf("starting user routine mode sleep handler...")

	for {
		select {
		// Channel to recieve any events.
		case got := <-in:
			d, _ := got.Data.(string)
			if got.Object == "mode_sleep" && d == "on" {
				m.ExecObject("living_light", "Off")
				m.ExecObject("dining_light", "Off")
			}
			// Run other code in default.
			//default:
			//m.Objects["ipaddress"].Execute("chil")
		}
	}
}
