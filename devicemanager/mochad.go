package devicemanager

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/deepakkamesh/viki/devicemanager/device"
	"github.com/deepakkamesh/viki/objectmanager"
	"github.com/golang/glog"
)

// Unique Device Id. Usually  same as device name.
const Device_MOCHAD DeviceId = "mochad"

type mochad struct {
	in     chan DeviceData
	quit   chan struct{}
	err    chan error
	out    chan DeviceData
	om     *objectmanager.ObjectManager
	conn   net.Conn
	ipPort string
}

// NewDeviceMochad returns a new and initialized mochad object.
func (m *DeviceSettings) NewDeviceMochad(out chan DeviceData, err chan error, om *objectmanager.ObjectManager) (DeviceId, device.Device) {
	return Device_MOCHAD, &mochad{
		in:   make(chan DeviceData, 10),
		quit: make(chan struct{}),
		err:  err,
		out:  out,
		om:   om,
	}
}

// Not implemented.
func (m *mochad) execute(data interface{}, object string) error {
	// Assert the command data depending on device.
	d, _ := data.(string)
	glog.Infof("mochad: executing %s on %s", d, object)
	return nil
}

// read will poll the mochad connection from any data and return one line of data.
func (m *mochad) runMochadPoll() {
	var err error
	re := regexp.MustCompile("(HouseUnit|Addr): (.+) Func: (.+)\n$")
	for {
		if m.conn == nil {
			if m.conn, err = net.DialTimeout("tcp", m.ipPort, time.Duration(5)*time.Second); err != nil {
				m.err <- fmt.Errorf("unable to connect to mochad %s", err)
				time.Sleep(60 * time.Second) // Sleep do we dont keep retrying too often.
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

			// Decode mochad state.
			state := strings.Trim(matches[3], " ")
			if strings.Contains(state, "normal") {
				state = "Closed"
			} else if strings.Contains(state, "alert") {
				state = "Open"
			}

			m.out <- DeviceData{
				DeviceId: Device_MOCHAD,
				Address:  strings.Trim(matches[2], " "),
				Data:     state,
			}
		}
	}
}

// On turns off the device.
// Not implemented.
func (m *mochad) On() {
	glog.Infof("Turn off")
}

// Off turns off the device.
// Not implemented.
func (m *mochad) Off() {
	glog.Infof("Turn off")
}

// Start initiates the device.
func (m *mochad) Start() error {
	glog.Infof("starting device mochad...")
	m.ipPort = flag.Lookup("mochad_ipport").Value.String()

	go m.run()
	go m.runMochadPoll()
	return nil
}

// Execute queues up the requested command to the channel.
func (m *mochad) Execute(action interface{}, object string) {
	m.in <- DeviceData{
		Data:    action,
		Address: object,
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
			if err := m.execute(in.Data, in.Address); err != nil {
				m.err <- err
			}

		case <-m.quit:
			return
		}
	}
}
