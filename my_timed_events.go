package viki

import (
	"log"
	"time"
	"viki/devicemanager"
)

func (m *Viki) timedEvents(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine timedEvents...")
	tick := time.NewTicker(1 * time.Second)
	for {
		select {
		// Channel to recieve any events.
		case <-c:
			//d, _ := got.Data.(string)
			//log.Printf("Got data from %s %s\n", got.Object, d)

		// Check for time changes.
		case <-tick.C:
			hhmm := time.Now().Format("1504")
			lrState, _ := m.Objects["living_light"].State.(string)
			drState, _ := m.Objects["dining_light"].State.(string)

			if hhmm == "1700" && lrState != "On" && drState != "On" {
				m.Objects["living_light"].Execute("On")
				m.Objects["dining_light"].Execute("On")
				log.Printf("turning on living and dining room lights")
			}

			// Run other code in default.
			//default:
		}
	}
}
