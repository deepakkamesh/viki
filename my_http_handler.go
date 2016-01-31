package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
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
				if err := m.ExecObject(d[0], state); err != nil {
					log.Printf("recieved unknown object %s", d[0])
					continue
				}
				m.ExecObject("speaker", "Executing command")
			}
			// Run other code in default.
			//default:
		}
	}
}
