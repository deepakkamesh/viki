package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

// Reference Format: Mon Jan 2 15:04:05 -0700 MST 2006
func (m *Viki) timedEvents(in chan devicemanager.DeviceData) {

	log.Printf("starting user routine timedEvents...")
	t1700 := NewReminder("1700", "1504") // Ping every 5pm.

	for {
		select {
		// Channel to recieve any events.
		case got := <-in:
			d, _ := got.Data.(string)
			log.Printf("Got data from %s %s\n", got.Object, d)

		// At 5pm.
		case <-t1700.C:
			m.ExecObject("living_light", "On")
			m.ExecObject("dining_light", "On")
			log.Printf("turning on living and dining room lights")

		}

	}

	// Run other code in default.
	//default:
}
