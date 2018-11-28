package viki

import (
	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/golang/glog"
)

// Reference Format: Mon Jan 2 15:04:05 -0700 MST 2006
func (m *Viki) MyTimedEvents(in chan devicemanager.DeviceData) {

	glog.Infof("Starting user routine MyTimedEvents...")
	defer glog.Infof("Shutting down user routine MyTimedEvents")

	t1900 := NewReminder("1900", "1504") // Ping every 5pm.
	t1700 := NewReminder("1700", "1504") // Ping every 5pm.
	t2200 := NewReminder("2200", "1504") // Ping every 10pm.
	t2000 := NewReminder("2000", "1504") // Ping every 8pm.
	t0500 := NewReminder("0500", "1504") // Ping every 5am.

	/* 	t0001 := NewReminder("0001", "1504") // Ping every 12:01am.
	var s sunrise.Sunrise
		lat := flag.Lookup("lat").Value.(flag.Getter).Get().(float64)
		long := flag.Lookup("long").Value.(flag.Getter).Get().(float64)
		s.Around(lat, long, time.Now())
		sunrise := s.Sunrise()
	*/
	_ = t1900

	for {
		select {
		// Channel to recieve any events.
		case <-in:

		// Turn off lights in the morning.
		case <-t0500.C:
			m.Do("living light", "Off")
			m.Do("dining light", "Off")
			m.Do("patio light", "Off")
			m.Do("tv light", "Off")
			glog.Infof("Turning off all lights")

		case <-t1700.C:
			m.Do("living light", "On")
			m.Do("dining light", "On")
			m.Do("patio light", "On")
			m.Do("tv light", "Off")
			glog.Infof("Turning on lights in evening")

		case <-t2000.C:
			if m.getModeState("mode vacation") == "Off" {
				m.Do("bedroom light", "On")
			}

		case <-t2200.C:
			m.Do("living light", "Off")
			if m.getModeState("mode vacation") == "Off" {
				m.Do("patio light", "Off")
				continue
			}
			m.Do("dining light", "Off")
		}
	}
}
