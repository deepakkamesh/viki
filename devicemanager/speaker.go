/* Speak uses festival to synthesis speech */
package devicemanager

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/deepakkamesh/viki/devicemanager/device"
	"github.com/deepakkamesh/viki/objectmanager"
	"github.com/golang/glog"
)

// Unique Device ID.
const Device_SPEAKER DeviceId = "speaker"

type speaker struct {
	in     chan DeviceData
	quit   chan struct{}
	err    chan error
	out    chan DeviceData
	om     *objectmanager.ObjectManager
	ipPort string
	conn   net.Conn
}

// NewDeviceSpeaker returns a new and initialized speaker.
func (m *DeviceSettings) NewDeviceSpeaker(out chan DeviceData, err chan error, om *objectmanager.ObjectManager) (DeviceId, device.Device) {
	return Device_SPEAKER, &speaker{
		in:   make(chan DeviceData, 10),
		quit: make(chan struct{}),
		err:  err,
		out:  out,
		om:   om,
	}
}

// speakFestival speaks the data on festival.
func (m *speaker) speakFestival(data interface{}) error {
	var err error
	text, _ := data.(string)
	glog.Infof("speaking  %s", text)

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
	glog.Infof("starting device speak...")
	flag := flag.Lookup("festival_ipport")
	m.ipPort = flag.Value.String()

	go m.run()
	return nil
}

// Execute queues up the requested command to the channel.
func (m *speaker) Execute(action interface{}, object string) {
	m.in <- DeviceData{
		Data:    action,
		Address: object,
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
			switch in.Address {
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
