/* Speak uses festival to synthesis speech */
package devicemanager

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

type Speaker struct {
	deviceNumber DeviceNumber
	cmd          chan DeviceData
	quit         chan struct{}
	err          chan error
	data         chan DeviceData
	ipPort       string
	conn         net.Conn
}

func (m *Speaker) speakFestival(data interface{}) error {
	var errW, err error
	r := 3

	text, _ := data.(string)
	log.Printf("speaking  %s", text)

	// If conn is not initialized, attempt to connect.
	if m.conn == nil {
		if m.conn, err = net.DialTimeout("tcp", m.ipPort, time.Duration(2)*time.Second); err != nil {
			return fmt.Errorf("unable to connect to %s:%s", m.ipPort, err)
		}
	}

	// Try to write to socket, if not try to connect and write again.
	for _, errW = fmt.Fprintf(m.conn, "(tts_text \"%s\" nil)", text); errW != nil && r > 0; r -= 1 {
		if m.conn, err = net.DialTimeout("tcp", m.ipPort, time.Duration(2)*time.Second); err != nil {
			return fmt.Errorf("unable to connect to %s:%s", m.ipPort, err)
		}
	}
	if errW != nil {
		return fmt.Errorf("unable to write to socket %s:%s", m.ipPort, err)
	}
	return nil
}

// On is not implemented.
func (m *Speaker) On() {
}

// Off is not implemented.
func (m *Speaker) Off() {
}

// Start initiates the device.
func (m *Speaker) Start() error {
	log.Printf("starting device speak...")
	flag := flag.Lookup("festival_ipport")
	m.ipPort = flag.Value.String()

	go m.run()
	return nil
}

// Execute queues up the requested command to the channel.
func (m *Speaker) Execute(action interface{}, object string) {
	m.cmd <- DeviceData{
		Data:   action,
		Object: object,
	}
}

// Shutdown terminates the device processing.
func (m *Speaker) Shutdown() {
	m.quit <- struct{}{}
}

// run is the main processing loop for the device driver.
func (m *Speaker) run() {
	for {
		select {
		case cmd := <-m.cmd:
			switch cmd.Object {
			case "festival":
				if err := m.speakFestival(cmd.Data); err != nil {
					m.err <- err
				}
			}
		case <-m.quit:
			return
		}
	}
}
