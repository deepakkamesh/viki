/*
Try to use the data coming from channel as the trigger and use the objects as conditions
to do some action.
This code is a template for writing custom actions
*/
package viki

import (
	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/golang/glog"
)

// User code starts with "my"
// Viki uses reflection to run any usercode starting with my*.
// Anything else is ignored.
func (m *Viki) MyUserCode(c chan devicemanager.DeviceData) {

	glog.Infof("Starting user routine userCode...")
	defer glog.Infof("Shutting down user routine userCode...")

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
