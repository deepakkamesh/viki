// +build ignore

/*
Try to use the data coming from channel as the trigger and use the objects as conditions
to do some action.
This code is a template for writing custom actions
*/
package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

// User code starts with "my"
// Viki uses reflection to run any usercode starting with my*.
// Anything else is ignored.
func (m *Viki) myUserCode(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine userCode...")

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			_ = got
			//d, _ := got.Data.(string)
			//log.Printf("Got data from %s %s\n", got.Object, d)

			// Run other code in default.
			//default:

			//m.Objects["ipaddress"].Execute("chil")
		}
	}
}
