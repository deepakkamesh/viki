/* modeSleep will turn off the lights and if there is any motion within the living
room, turn on the living room lights. If there is any external motion, trigger alarm
*/
package viki

import (
	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/golang/glog"
)

func (m *Viki) MyModeMovie(in chan devicemanager.DeviceData) {

	glog.Infof("Starting user routine MyModeMovie...")
	glog.Infof("Shutting down user routine MyModeMovie")

	for {
		select {
		case got := <-in:
			name, _ := m.ObjectManager.GetObjectByAddress(got.Address)
			d := m.getMochadState(name)
			if got.Address == "mode_movie" {
				if d == "On" {
					m.Do("living light", "Off")
					m.Do("dining light", "Off")
					m.Do("tv light", "On")
				}
				if d == "Off" {
					m.Do("living light", "On")
					m.Do("dining light", "On")
					m.Do("tv light", "Off")
				}
			}
		}
	}
}
