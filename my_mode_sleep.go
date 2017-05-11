package viki

import (
	"fmt"
	"log"
	"time"

	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/mailgun/mailgun-go"
)

/* MyModeSleep will turn off the lights and if there is any motion within the living
room, turn on the living room lights. If there is any external motion, trigger alarm
*/
func (m *Viki) MyModeSleep(in chan devicemanager.DeviceData) {

	log.Printf("starting user routine mode sleep handler...")
	t0600 := NewReminder("0600", "1504") // Ping every 5am.
	mg := mailgun.NewMailgun("sandboxf139420cc83d4d3a8c3cf5dfc9b06b42.mailgun.org", "key-6ceddfaf05c0d237076a19abe2afef5d", "pubkey-ce009cba9207ec56ae09ac45b9607c2f")
	for {
		select {
		// Channel to recieve any events.
		case got := <-in:
			name, obj := m.ObjectManager.GetObjectByAddress(got.Object)
			st := m.getMochadState(name)
			if got.Object == "mode_sleep" && st == "On" {
				m.execObject("living light", "Off")
				m.execObject("dining light", "Off")
				m.execObject("tv light", "Off")
				continue
			}

			if m.getModeState("mode sleep") == "On" {
				// Setup some alerting when sleeping.
				if st == "Open" && obj.CheckTag("door") {
					msg := fmt.Sprintf("%s Open", name)
					quickMail("deepak.kamesh@gmail.com", msg, mg)
					quickMail("6024050044@tmomail.net", msg, mg)
					// for a bit.
					m.execObject("living light", "On")
					m.execObject("dining light", "On")
					m.execObject("bedroom light", "On")
					m.execObject("buzzer", "On")
					time.AfterFunc(15*time.Minute, func() {
						m.execObject("living light", "Off")
						m.execObject("dining light", "Off")
						m.execObject("bedroom light", "Off")
						m.execObject("buzzer", "Off")
					})
					continue
				}

				// Setup automatic turn on of lights.
				if st == "On" && obj.CheckTag("indoor_motion") {
					// Turn on the living room light for a bit.
					m.execObject("living light", "On")
					time.AfterFunc(3*time.Minute, func() {
						m.execObject("living light", "Off")
					})
					continue
				}
			}

		// Turn off mode sleep at 5am.
		case <-t0600.C:
			m.execObject("mode sleep", "Off")

		}
	}
}
