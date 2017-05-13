/*
Try to use the data coming from channel as the trigger and use the objects as conditions
to do some action.
This code is a template for writing custom actions
*/
package viki

import (
	"encoding/json"

	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/golang/glog"
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

	glog.Infof("starting user routine MyArtikHandler...")
	defer glog.Infof("starting user routine MyArtikHandler")

	for {
		select {
		case got := <-c:
			if got.Address != "artik" {
				continue
			}
			d, _ := got.Data.([]byte)

			aa := ArtikAction{}
			if err := json.Unmarshal(d, &aa); err != nil {
				glog.Errorf("Error unmarshalling: %v", err)
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
