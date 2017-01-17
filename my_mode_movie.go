/* modeSleep will turn off the lights and if there is any motion within the living
room, turn on the living room lights. If there is any external motion, trigger alarm
*/
package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

func (m *Viki) MyModeMovie(in chan devicemanager.DeviceData) {

	log.Printf("starting user routine mode sleep handler...")

	for {
		select {
		// Channel to recieve any events.
		case got := <-in:
			d, _ := got.Data.(string)
			if got.Object == "mode_movie" {
				if d == "On" {
					m.execObject("living light", "Off")
					m.execObject("dining light", "Off")
					m.execObject("tv light", "On")
				}
				if d == "Off" {
					m.execObject("living light", "On")
					m.execObject("dining light", "On")
					m.execObject("tv light", "Off")
				}
			}
		}
	}
}
