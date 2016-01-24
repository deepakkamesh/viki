/* Package viki is a intelligent extensible home automation framework.
 */
package viki

import (
	"fmt"
	"log"
	"time"
	"viki/devicemanager"
)

type UserCode struct {
	f    func(chan devicemanager.DeviceData)
	data chan devicemanager.DeviceData
}

type Viki struct {
	Version       int
	Objects       map[string]*Object
	DeviceManager *devicemanager.DeviceSettings
	UserCodes     []*UserCode
}

// ReadConfig reads the configuration from configuration file.
func (m *Viki) readConfig(file string) error {

	// TODO(dkg): Read configuration from file.
	m.Objects["living_room"] = InitObject("living_room", "C4", devicemanager.Device_MISTERHOUSE, m.DeviceManager)
	m.Objects["dining_room"] = InitObject("dining_room", "M1", devicemanager.Device_MISTERHOUSE, m.DeviceManager)
	m.Objects["ipaddress"] = InitObject("ipaddress", "ipaddress", devicemanager.Device_HTTPHANDLER, m.DeviceManager)

	return nil
}

func New() *Viki {
	return &Viki{
		Version: 1,
		Objects: make(map[string]*Object),
	}
}

func (m *Viki) Init() error {

	// Initialize device manager.
	m.DeviceManager = devicemanager.New()

	// Initiatilze User Code.
	m.UserCodes = []*UserCode{
		&UserCode{
			f:    m.doSomething,
			data: make(chan devicemanager.DeviceData),
		},
	}

	// Read configuration.
	if err := m.readConfig("fixme"); err != nil {
		return fmt.Errorf("unable to open configuration file %s", err)
	}
	return nil
}

func (m *Viki) Run() {
	// Start Device Manager.
	m.DeviceManager.StartDeviceManager()

	// Start User Code.
	for _, userCode := range m.UserCodes {
		go userCode.f(userCode.data)
	}

	// Run the main processing loop.
	for {
		select {
		case data := <-m.DeviceManager.Data:
			//	 Send recieved data to all user code channels.
			for _, userCode := range m.UserCodes {
				userCode.data <- data
			}

		case err := <-m.DeviceManager.Err:
			log.Printf("device manager error %s", err)
		}
	}
}

// User Code.

func (m *Viki) doSomething(c chan devicemanager.DeviceData) {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case devData := <-c:
			d, _ := devData.Data.(string)
			fmt.Printf("GOt data from %s %s", data.Object, d)
		case <-t.C:
			m.Objects["ipaddress"].Execute("chil")
		default:
			//m.Objects["ipaddress"].Execute("chil")
		}
	}
}
