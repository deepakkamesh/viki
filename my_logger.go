package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

func (m *Viki) logger(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine logger...")

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			d, _ := got.Data.(string)
			name := m.GetObjectName(got.Object)
			if name != "" {
				log.Printf("Got data from %s %s\n", name, d)
			} else {
				log.Printf("Got data from %s %s\n", got.Object, d)
			}
			// Run other code in default.
			//default:
			//m.Objects["ipaddress"].Execute("chil")
		}
	}
}
