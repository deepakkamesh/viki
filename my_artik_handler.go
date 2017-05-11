/*
Try to use the data coming from channel as the trigger and use the objects as conditions
to do some action.
This code is a template for writing custom actions
*/
package viki

import (
	"encoding/json"
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

type ArtikAction struct {
	Type string `json:"type"`
	Cts  int64  `json:"cts"`
	Ts   int64  `json:"ts"`
	Mid  string `json:"mid"`
	Sdid string `json:"sdid"`
	Ddid string `json:"ddid"`
	Data struct {
		Actions []struct {
			Name string `json:"name"`
		} `json:"actions"`
	} `json:"data"`
	Ddtid string `json:"ddtid"`
	UID   string `json:"uid"`
	Boid  string `json:"boid"`
	Mv    int    `json:"mv"`
}

func (m *Viki) MyArtikHandler(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine artikHandler...")

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			if got.Object != "artik" {
				continue
			}
			d, _ := got.Data.([]byte)

			aa := ArtikAction{}
			if err := json.Unmarshal(d, &aa); err != nil {
				log.Printf("error unmarshalling %v", err)
				continue
			}

			ddid := aa.Ddid
			for _, o := range m.ObjectManager.Objects {
				if _, ok := o.GetTag(ddid); !ok {
					continue
				}
				for _, cmd := range aa.Data.Actions {
					switch cmd.Name {
					case "setOn":
						o.Execute("On")
					case "setOff":
						o.Execute("Off")
					}
				}
			}
		}
	}
}
