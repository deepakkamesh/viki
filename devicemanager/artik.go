/* Example stub code for a device driver.
Steps to add a new device driver.
1. Make a copy of this file and update.
2. Add the device to devices.go
*/
package devicemanager

import (
	"fmt"
	"log"
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
		inC:  make(chan []byte), //Input channel for actions
	}
}

func (m *artik) execute(data interface{}, object string) error {

	// Assert the command data depending on device.
	d, _ := data.(string)
	log.Printf("artik: executing %d on %s", d, object)
	return nil
}

// On turns off the device.
func (m *artik) On() {
	log.Printf("Turn off")
}

// Off turns off the device.
func (m *artik) Off() {
	log.Printf("Turn off")
}

// Start initiates the device.
func (m *artik) Start() error {
	log.Printf("starting device artik...")
	// Set any required parameters using flag.
	m.register()
	go m.receiver()
	go m.run()
	return nil
}

// Execute queues up the requested command to the channel.
func (m *artik) Execute(action interface{}, object string) {
	m.in <- DeviceData{
		Data:   action,
		Object: object,
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
			if err := m.execute(in.Data, in.Object); err != nil {
				m.err <- err
			}
		// Actions from Artik Cloud.
		case msg := <-m.inC:
			m.out <- DeviceData{
				DeviceId: Device_ARTIK,
				Data:     msg,
				Object:   "artik",
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
			log.Println("read:", err)
			return
		}
		glog.V(2).Infof("recv: %s", message)
		m.inC <- message
	}
}

func (m *artik) register() {
	u := url.URL{Scheme: "wss", Host: "api.artik.cloud", Path: "/v1.1/websocket"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	m.ws = c
	// register device.

	//	regMsg := `{"type":"register", "sdid":"d2381f135c6441a2a895f1341c969ec6", "Authorization":"bearer 0d434ab87f4c4586898520d96afa3181", "cid":"1234546"}`

	for _, o := range m.om.Objects {
		deviceID, ok1 := o.GetTag("ddid")
		deviceToken, ok2 := o.GetTag("dtid")
		if !ok1 || !ok2 {
			continue
		}
		ddid := strings.Split(deviceID, "=")[1]
		dtid := strings.Split(deviceToken, "=")[1]

		regMsg := fmt.Sprintf("{\"type\":\"register\", \"sdid\":\"%s\", \"Authorization\":\"bearer %s\"}", ddid, dtid)
		log.Printf("Register %v", regMsg)
		c.WriteMessage(websocket.TextMessage, []byte(regMsg))
	}
}
