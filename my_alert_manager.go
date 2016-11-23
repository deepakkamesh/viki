package viki

import (
	"fmt"
	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/mailgun/mailgun-go"
	"log"
	"time"
)

func (m *Viki) MyAlertManager(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine Alert Manager...")
	// TODO: Load from flags.
	mg := mailgun.NewMailgun("sandboxf139420cc83d4d3a8c3cf5dfc9b06b42.mailgun.org", "key-6ceddfaf05c0d237076a19abe2afef5d", "pubkey-ce009cba9207ec56ae09ac45b9607c2f")

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			name, obj := m.getObject(got.Object)

			// Alerts when we are not at home.
			if m.getModeState("mode vacation") == "On" {
				st := m.getMochadState(name)
				// Motion inside.
				//if st == "On" && (name == "dining_ms1" || name == "living_ms1" || name == "bedroom_ms1") {
				if st == "On" && obj.checkTag("indoor_motion") {
					msg := fmt.Sprintf("Detected motion in %s", name)
					quickMail("deepak.kamesh@gmail.com", msg, mg)
					quickMail("6024050044@tmomail.net", msg, mg)
					// Turn on the living room light for a bit.
					m.execObject("living light", "On")
					m.execObject("dining light", "On")
					m.execObject("buzzer", "On")
					time.AfterFunc(3*time.Minute, func() {
						m.execObject("living light", "Off")
						m.execObject("dining light", "Off")
						m.execObject("buzzer", "Off")
					})

					continue
				}
				// Doors opened.
				//if st == "Open" && (name == "backyard door" || name == "garage door" || name == "front door") {
				if st == "Open" && obj.checkTag("door") {
					msg := fmt.Sprintf("%s Open", name)
					quickMail("deepak.kamesh@gmail.com", msg, mg)
					quickMail("6024050044@tmomail.net", msg, mg)
					continue
					// for a bit.
					m.execObject("living light", "On")
					m.execObject("dining light", "On")
					m.execObject("buzzer", "On")
					time.AfterFunc(3*time.Minute, func() {
						m.execObject("living light", "Off")
						m.execObject("dining light", "Off")
						m.execObject("buzzer", "Off")
					})

					continue
				}

			}
			/*TODO - Throttle warning to no more than once every 5s.
			// if motion sensor backyard and door is not open.
			if st == "On" && name == "backyard_ms1" && m.getMochadState("backyard door") != "Open" {
				m.execObject("speaker", "Warning backyard motion sensor activated ")
				continue
			}
			// if motion sensor garage and door is not open.
			if st == "On" && name == "garage_ms1" && m.getMochadState("garage door") != "Open" {
				m.execObject("speaker", "Warning garage motion sensor activated ")
				continue
			}
			*/

		}
	}
}

func quickMail(recipient string, msg string, mg mailgun.Mailgun) {
	message := mailgun.NewMessage("home@fulton-ave", "Alert!", msg, recipient)
	msg, id, err := mg.Send(message)
	if err != nil {
		log.Printf("Could not send message: %v, ID %v, %+v", err, id, msg)
	}

}
