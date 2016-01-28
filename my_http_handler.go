package viki

import (
	"log"
	"viki/devicemanager"
)

func (m *Viki) httpHandler(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine httphandler...")

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			if got.Object == "http" {
				d, _ := got.Data.([]string)
				state := sanitizeState(d[1])
				log.Printf("Got data from %s %s\n", got.Object, d)
				if obj, ok := m.Objects[d[0]]; ok {
					obj.Execute(state)
					m.Objects["speaker"].Execute("Executing command")
					continue
				}
				log.Printf("recieved unknown object %s", d[0])
			}
			// Run other code in default.
			//default:
		}
	}
}
