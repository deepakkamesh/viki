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
		case got := <-c:
			d, _ := got.Data.(string)
			log.Printf("Got data from %s %s\n", got.Object, d)

		// Check for time changes.
		case <-tick.C:
			hhmm := time.Now().Format("1504")
			lrState, _ := m.Objects["living_light"].State.(string)
			if hhmm == "2038" && lrState != "On" {
				m.Objects["living_light"].Execute("On")
			}

		// Run other code in default.
		default:
			//m.Objects["ipaddress"].Execute("chil")
		}
	}
}
