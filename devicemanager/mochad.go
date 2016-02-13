package devicemanager

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

// Unique Device Id. Usually  same as device name.
const Device_MOCHAD DeviceId = "mochad"

type mochad struct {
	deviceId DeviceId
	in       chan DeviceData
	quit     chan struct{}
	err      chan error
	out      chan DeviceData
	conn     net.Conn
	ipPort   string
}

// Not implemented.
func (m *mochad) execute(data interface{}, object string) error {
	// Assert the command data depending on device.
	d, _ := data.(string)
	log.Printf("mochad: executing %d on %s", d, object)
	return nil
}

// read will poll the mochad connection from any data and return one line of data.
func (m *mochad) runMochadPoll() {
	var err error
	re := regexp.MustCompile("(HouseUnit|Addr): (.+) Func: Contact_(.+)_(min|max)_DS10A\n$")
	for {
		if m.conn == nil {
			if m.conn, err = net.DialTimeout("tcp", m.ipPort, time.Duration(5)*time.Second); err != nil {
				m.err <- fmt.Errorf("unable to connect to mochad %s", err)
				time.Sleep(5 * time.Second) // Sleep do we dont keep retrying too often.
				continue
			}
		}
		buf, err := bufio.NewReader(m.conn).ReadString('\n')
		if err != nil {
			// If error lets close connect and reconnect.
			m.err <- fmt.Errorf("unable to read mochad %s", err)
			m.conn.Close()
			m.conn = nil
			continue
		}
		matches := re.FindStringSubmatch(buf)
		if matches != nil {
			m.out <- DeviceData{
				DeviceId: m.deviceId,
				Object:   strings.Trim(matches[2], " "),
				Data:     strings.Trim(matches[3], " "),
			}
		}
	}
}

// On turns off the device.
// Not implemented.
func (m *mochad) On() {
	log.Printf("Turn off")
}

// Off turns off the device.
// Not implemented.
func (m *mochad) Off() {
	log.Printf("Turn off")
}

// Start initiates the device.
func (m *mochad) Start() error {
	log.Printf("starting device mochad...")
	m.ipPort = flag.Lookup("mochad_ipport").Value.String()

	go m.run()
	go m.runMochadPoll()
	return nil
}

// Execute queues up the requested command to the channel.
func (m *mochad) Execute(action interface{}, object string) {
	m.in <- DeviceData{
		Data:   action,
		Object: object,
	}
}

// Shutdown terminates the device processing.
func (m *mochad) Shutdown() {
	m.quit <- struct{}{}
}

// run is the main processing loop for the device driver.
func (m *mochad) run() {
	for {
		select {
		case in := <-m.in:
			if err := m.execute(in.Data, in.Object); err != nil {
				m.err <- err
			}

		case <-m.quit:
			return
		}
	}
}
