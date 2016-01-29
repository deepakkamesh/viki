package viki

import "time"

type reminder struct {
	C    <-chan time.Time
	quit chan struct{}
}

// NewReminder sends a reminder every at time specified.
// at is interpreted using format.
// For eg. at = 2030 format = 1504, ping every 8.30pm.
// Reference Format: Mon Jan 2 15:04:05 -0700 MST 2006
func NewReminder(at string, format string) *reminder {

	tick := time.NewTicker(250 * time.Millisecond)
	ping := make(chan time.Time)
	quit := make(chan struct{})

	go func() {
		fl := true
		for {
			select {
			case <-tick.C:
				if time.Now().Format(format) == at {
					if fl {
						ping <- time.Now()
						fl = false
					}
					continue
				}
				fl = true
			case <-quit:
				return
			}
		}
	}()

	return &reminder{
		C:    ping,
		quit: quit,
	}
}

func (m *reminder) Stop() {
	m.quit <- struct{}{}
}
