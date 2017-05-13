package viki

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/golang/glog"
)

func (m *Viki) MyAlertManager(c chan devicemanager.DeviceData) {

	glog.Infof("Starting user routine MyAlertManager...")
	defer glog.Infof("Shutting down user routine MyAlertManage...")

	fgRecipients := flag.Lookup("email_alert_list")

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			name, obj := m.ObjectManager.GetObjectByAddress(got.Address)

			// Alerts when we are not at home.
			if m.getModeState("mode vacation") == "On" {
				st := m.getMochadState(name)
				// Motion inside.
				if st == "On" && obj.CheckTag("indoor_motion") {
					msg := fmt.Sprintf("Detected motion in %s", name)
					if fgRecipients != nil {
						if err := m.quickMail(strings.Split(fgRecipients.Value.String(), ","), msg); err != nil {
							glog.Errorf("failed to send email %v", err)
						}
					}

					m.Do("living light", "On")
					m.Do("dining light", "On")
					m.Do("buzzer", "On")
					time.AfterFunc(3*time.Minute, func() {
						m.Do("living light", "Off")
						m.Do("dining light", "Off")
						m.Do("buzzer", "Off")
					})
					continue
				}

				// Doors opened.
				if st == "Open" && obj.CheckTag("door") {
					msg := fmt.Sprintf("%s Open", name)
					if fgRecipients != nil {
						if err := m.quickMail(strings.Split(fgRecipients.Value.String(), ","), msg); err != nil {
							glog.Errorf("failed to send email %v", err)
						}
					}

					m.Do("living light", "On")
					m.Do("dining light", "On")
					m.Do("buzzer", "On")
					time.AfterFunc(3*time.Minute, func() {
						m.Do("living light", "Off")
						m.Do("dining light", "Off")
						m.Do("buzzer", "Off")
					})
					continue
				}
			}
		}
	}
}
