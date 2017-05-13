/* Example stub code for a device driver.
Steps to add a new device driver.
1. Make a copy of this file and update.
2. Add the device to devices.go
*/
package devicemanager

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/deepakkamesh/viki/devicemanager/device"
	"github.com/deepakkamesh/viki/objectmanager"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

// Unique Device Id. Usually  same as device name.
const Device_ARTIK DeviceId = "artik"

type artik struct {
	in   chan DeviceData
	quit chan struct{}
	err  chan error
	out  chan DeviceData
	om   *objectmanager.ObjectManager
	inC  chan []byte
	ws   *websocket.Conn
}

func (m *DeviceSettings) NewDeviceArtik(out chan DeviceData, err chan error, om *objectmanager.ObjectManager) (DeviceId, device.Device) {
	return Device_ARTIK, &artik{
		in:   make(chan DeviceData, 10), // Input channel, typically buffered.
		quit: make(chan struct{}),       // Quit.
		err:  err,                       // Common error channel.
		out:  out,                       // Channel to send out data.
		om:   om,
		inC:  make(chan []byte, 10), //Input channel for actions
	}
}

func (m *artik) execute(data interface{}, object string) error {

	// Assert the command data depending on device.
	d, _ := data.(string)
	glog.Infof("artik: executing %d on %s", d, object)
	return nil
}

// On turns off the device.
func (m *artik) On() {
	glog.Infof("Turn off")
}

// Off turns off the device.
func (m *artik) Off() {
	glog.Infof("Turn off")
}

// Start initiates the device.
func (m *artik) Start() error {
	glog.Infof("Starting device artik...")
	if err := m.register(); err != nil {
		fmt.Errorf("registration failed %v", err)
	}
	go m.receiver()
	go m.run()
	return nil
}

// Execute queues up the requested command to the channel.
func (m *artik) Execute(action interface{}, object string) {
	m.in <- DeviceData{
		Data:    action,
		Address: object,
	}
}

// Shutdown terminates the device processing.
func (m *artik) Shutdown() {
	m.ws.Close()
	m.quit <- struct{}{}
}

// run is the main processing loop for the device driver.
func (m *artik) run() {
	for {
		select {
		// Send sensor data to Artik Cloud.
		case in := <-m.in:
			if err := m.execute(in.Data, in.Address); err != nil {
				m.err <- err
			}
		// Actions from Artik Cloud.
		case msg := <-m.inC:
			m.out <- DeviceData{
				DeviceId: Device_ARTIK,
				Data:     msg,
				Address:  "artik",
			}

		case <-m.quit:
			return
		}
	}
}

func (m *artik) receiver() {
	for {
		_, message, err := m.ws.ReadMessage()
		if err != nil {
			glog.Errorf("Read failed from websock: %v", err)
			continue
		}
		glog.V(2).Infof("Recieved on websocket: %s", message)
		m.inC <- message
	}
}

func (m *artik) register() error {
	u := url.URL{Scheme: "wss", Host: "api.artik.cloud", Path: "/v1.1/websocket"}
	glog.V(1).Infof("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to dial Artik: %v", err)
	}
	m.ws = c

	// register device.
	for _, o := range m.om.Objects {
		deviceID, ok1 := o.GetTag("ddid")
		deviceToken, ok2 := o.GetTag("dtid")
		if !(ok1 && ok2) {
			continue
		}
		ddid := strings.Split(deviceID, "=")[1]
		dtid := strings.Split(deviceToken, "=")[1]

		regMsg := fmt.Sprintf("{\"type\":\"register\", \"sdid\":\"%s\", \"Authorization\":\"bearer %s\"}", ddid, dtid)
		glog.V(1).Infof("Registering device %s Msg:%v", o.Name, regMsg)
		if err := c.WriteMessage(websocket.TextMessage, []byte(regMsg)); err != nil {
			return fmt.Errorf("failed to write message: %v", err)
		}
	}
	return nil
}
