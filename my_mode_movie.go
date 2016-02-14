/* modeSleep will turn off the lights and if there is any motion within the living
room, turn on the living room lights. If there is any external motion, trigger alarm
*/
package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

func (m *Viki) modeMovie(in chan devicemanager.DeviceData) {

	log.Printf("starting user routine mode sleep handler...")

	for {
		select {
		// Channel to recieve any events.
		case got := <-in:
			d, _ := got.Data.(string)
			if got.Object == "mode_movie" {
				if d == "On" {
					m.ExecObject("living light", "Off")
					m.ExecObject("dining light", "Off")
					m.ExecObject("tv light", "On")
				}
				if d == "Off" {
					m.ExecObject("living light", "On")
					m.ExecObject("dining light", "On")
					m.ExecObject("tv light", "Off")
				}
			}
		}
	}
}
