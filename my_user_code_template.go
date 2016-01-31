/*
Try to use the data coming from channel as the trigger and use the objects as conditions
to do some action.
*/
package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

func (m *Viki) userCode(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine userCode...")

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			d, _ := got.Data.(string)
			log.Printf("Got data from %s %s\n", got.Object, d)

		// Run other code in default.
		default:
			//m.Objects["ipaddress"].Execute("chil")
		}
	}
}
