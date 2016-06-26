// +build ignore
package viki

import (
	"flag"
	"log"
	"os/exec"
	"time"

	"github.com/deepakkamesh/viki/devicemanager"
)

// MyDpmsControl detects motion and turns the display screen on or off
// on the external display by using dpms shell commands.
func (m *Viki) DEPRECATED_myDpmsControl(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine dpms control...")
	res := flag.Lookup("resource").Value.String()
	screenOn := false
	for {
		select {
		case got := <-c:
			d, _ := got.Data.(string)
			_, o := m.getObject(got.Object)

			// Got some motion.
			if o.checkTag("motion") && d == "On" && !screenOn {
				// Turn on screen.
				if err := exec.Command(res + "/dpmsoff.sh").Run(); err != nil {
					log.Printf("error running dpms off %s ", err)
					continue
				}
				screenOn = true
				time.AfterFunc(60*time.Minute, func() {
					// Turn off screen.
					if err := exec.Command(res + "/dpmson.sh").Run(); err != nil {
						log.Printf("error running dpms on %s", err)
						return
					}
					screenOn = false
				})
			}
		}
	}
}
