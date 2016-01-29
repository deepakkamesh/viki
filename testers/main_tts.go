package main

import (
	"fmt"
	"time"
)

type reminder struct {
	C    <-chan time.Time
	quit chan struct{}
}

func NewReminder(at string, format string) *reminder {

	tick := time.NewTicker(250 * time.Millisecond)
	ping := make(chan time.Time)
	quit := make(chan struct{})

	go func() {
		fl := true
		for {
			select {
			case <-tick.C:
				fmt.Printf("Time %s\n", time.Now().Format(format))
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

//Mon Jan 2 15:04:05 -0700 MST 2006
func main() {

	tim := NewReminder("01", "05")
	for {
		t := <-tim.C
		fmt.Println("NOW", t)
		tim.Stop()
	}
}
