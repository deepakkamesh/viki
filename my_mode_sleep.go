/* modeSleep will turn off the lights and if there is any motion within the living
room, turn on the living room lights. If there is any external motion, trigger alarm
*/
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
			if got.Object == "mode_sleep" && d == "On" {
				m.ExecObject("living light", "Off")
				m.ExecObject("dining light", "Off")
			}
		}
	}
}
