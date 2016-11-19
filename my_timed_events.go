package viki

import (
	"log"

	"github.com/deepakkamesh/viki/devicemanager"
)

// Reference Format: Mon Jan 2 15:04:05 -0700 MST 2006
func (m *Viki) MyTimedEvents(in chan devicemanager.DeviceData) {

	log.Printf("starting user routine timedEvents...")
	t1700 := NewReminder("1700", "1504") // Ping every 5pm.
	t1900 := NewReminder("1900", "1504") // Ping every 5pm.
	t2200 := NewReminder("2200", "1504") // Ping every 10pm.
	t2000 := NewReminder("2000", "1504") // Ping every 8pm.
	/* 	t0001 := NewReminder("0001", "1504") // Ping every 12:01am.
	var s sunrise.Sunrise
		lat := flag.Lookup("lat").Value.(flag.Getter).Get().(float64)
		long := flag.Lookup("long").Value.(flag.Getter).Get().(float64)
		s.Around(lat, long, time.Now())
		sunrise := s.Sunrise()
	*/
	for {
		select {
		// Channel to recieve any events.
		case <-in:

		// At 5pm.
		case <-t1700.C:
			m.execObject("living light", "On")
			m.execObject("dining light", "On")
			m.execObject("patio light", "On")
			m.execObject("tv light", "On")
			log.Printf("turning on evening lights")
		case <-t2000.C:
			m.execObject("bedroom light", "On")
		case <-t2200.C:
			m.execObject("patio light", "Off")
			m.execObject("living light", "Off")

		}
	}
}
