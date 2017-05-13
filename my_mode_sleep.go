package viki

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/golang/glog"
)

/* MyModeSleep will turn off the lights and if there is any motion within the living
room, turn on the living room lights. If there is any external motion, trigger alarm
*/
func (m *Viki) MyModeSleep(in chan devicemanager.DeviceData) {

	glog.Infof("Starting user routine MyModeSleep...")
	defer glog.Infof("Shutting down user routine mMyModeSlee")

	t0600 := NewReminder("0600", "1504") // Ping every 5am.

	fgRecipients := flag.Lookup("alert_email_list")

	for {
		select {
		case got := <-in:
			name, obj := m.ObjectManager.GetObjectByAddress(got.Address)
			st := m.getMochadState(name)
			if got.Address == "mode_sleep" && st == "On" {
				m.Do("living light", "Off")
				m.Do("dining light", "Off")
				m.Do("tv light", "Off")
				continue
			}

			if m.getModeState("mode sleep") == "On" {
				// Setup some alerting when sleeping.
				if st == "Open" && obj.CheckTag("door") {
					msg := fmt.Sprintf("%s Open", name)
					if fgRecipients != nil {
						if err := m.quickMail(strings.Split(fgRecipients.Value.String(), ","), msg); err != nil {
							glog.Errorf("failed to send email %v", err)
						}
					}
					m.Do("living light", "On")
					m.Do("dining light", "On")
					m.Do("bedroom light", "On")
					m.Do("buzzer", "On")
					time.AfterFunc(15*time.Minute, func() {
						m.Do("living light", "Off")
						m.Do("dining light", "Off")
						m.Do("bedroom light", "Off")
						m.Do("buzzer", "Off")
					})
					continue
				}

				// Setup automatic turn on of lights.
				if st == "On" && obj.CheckTag("indoor_motion") {
					// Turn on the living room light for a bit.
					m.Do("living light", "On")
					time.AfterFunc(3*time.Minute, func() {
						m.Do("living light", "Off")
					})
					continue
				}
			}

		// Turn off mode sleep at 5am.
		case <-t0600.C:
			m.Do("mode sleep", "Off")

		}
	}
}
