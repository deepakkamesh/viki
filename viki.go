/* Package viki is a intelligent extensible home automation framework.
 */
package viki

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

type Config struct {
	Objects []Object
}

// ReadConfig reads the configuration from configuration file.
func (m *Viki) readConfig(file string) error {
	config := Config{}
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		return fmt.Errorf("error parsing config file ", err)
	}
	for _, o := range config.Objects {
		m.Objects[o.Name] = InitObject(o.Name, o.Address, o.DevNo, m.DeviceManager)
	}

	log.Printf("%+v", m.Objects)

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

	// Initiatilze user code.
	m.UserCodes = []*UserCode{
		&UserCode{
			f:    m.timedEvents,
			data: make(chan devicemanager.DeviceData),
		},
	}

	// Read configuration.
	if err := m.readConfig("../config.objects"); err != nil {
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
		case got := <-m.DeviceManager.Data:
			name, err := m.GetNameOfObject(got.Object)
			if err != nil {
				log.Println(err)
				continue
			}
			// Set state.
			m.Objects[name].State = got.Data
			// Send recieved data to all user code channels.
			for _, userCode := range m.UserCodes {
				userCode.data <- got
			}

		case err := <-m.DeviceManager.Err:
			log.Printf("device manager error %s", err)
		}
	}
}

func (m *Viki) GetNameOfObject(address string) (string, error) {

	for k, v := range m.Objects {
		if v.Address == address {
			return k, nil
		}
	}
	return "", fmt.Errorf("object with address %s not found", address)
}
