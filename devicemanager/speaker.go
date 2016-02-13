/* Speak uses festival to synthesis speech */
package devicemanager

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

// Unique Device ID.
const Device_SPEAKER DeviceId = "speaker"

type speaker struct {
	deviceId DeviceId
	in       chan DeviceData
	quit     chan struct{}
	err      chan error
	out      chan DeviceData
	ipPort   string
	conn     net.Conn
}

func (m *speaker) speakFestival(data interface{}) error {
	var err error
	text, _ := data.(string)
	log.Printf("speaking  %s", text)

	for r := 3; r > 0; r -= 1 {
		// If conn is not initialized, attempt to connect.
		if m.conn == nil {
			if m.conn, err = net.DialTimeout("tcp", m.ipPort, time.Duration(2)*time.Second); err != nil {
				return fmt.Errorf("unable to dial festival %s", err)
			}
		}

		if _, err := fmt.Fprintf(m.conn, "(tts_text \"%s\" nil)", text); err == nil {
			return nil
		}
		// Try to write to socket, if not try to reconnect and write again.
		m.conn.Close()
		m.conn = nil
	}
	return fmt.Errorf("unable to speak on festival %s", err)
}

// On is not implemented.
func (m *speaker) On() {
}

// Off is not implemented.
func (m *speaker) Off() {
}

// Start initiates the device.
func (m *speaker) Start() error {
	log.Printf("starting device speak...")
	flag := flag.Lookup("festival_ipport")
	m.ipPort = flag.Value.String()

	go m.run()
	return nil
}

// Execute queues up the requested command to the channel.
func (m *speaker) Execute(action interface{}, object string) {
	m.in <- DeviceData{
		Data:   action,
		Object: object,
	}
}

// Shutdown terminates the device processing.
func (m *speaker) Shutdown() {
	m.quit <- struct{}{}
}

// run is the main processing loop for the device driver.
func (m *speaker) run() {
	for {
		select {
		case in := <-m.in:
			switch in.Object {
			case "festival":
				if err := m.speakFestival(in.Data); err != nil {
					m.err <- err
				}
			}
		case <-m.quit:
			return
		}
	}
}
