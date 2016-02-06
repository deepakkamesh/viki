/* Package viki is a intelligent extensible home automation framework.
 */
package viki

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/deepakkamesh/viki/devicemanager"
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
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	c := []string{}
	comment := regexp.MustCompile(`^#|^[\s]*$`) // Ignore comment or blank lines.
	for scanner.Scan() {
		line := scanner.Text()
		if comment.MatchString(line) {
			continue
		}
		c = strings.Split(line, ",")
		for i, _ := range c {
			c[i] = strings.Trim(c[i], " ")
		}

		// Get device if any.
		var (
			ok  bool
			dev devicemanager.Device
		)
		i := 2
		// Ignore device if device not specified or empty.
		if len(c)-1 >= i && len(c[i]) > 0 {
			if dev, ok = m.DeviceManager.Devices[devicemanager.DeviceId(c[i])]; !ok {
				return fmt.Errorf("invalid device \"%s\" specified", c[i])
			}
		}

		// Get tags if any.
		tags, i := []string{}, 3
		if len(c)-1 >= i {
			tags = strings.Split(c[i], "|")
			for j, _ := range tags {
				tags[j] = strings.Trim(tags[j], " ")
			}
		}
		m.Objects[c[1]] = InitObject(c[0], dev, tags)
	}
	return nil
}

func New() *Viki {
	return &Viki{
		Version: 1,
		Objects: make(map[string]*Object),
	}
}

func (m *Viki) Init(configFile string) error {

	// Initialize device manager.
	m.DeviceManager = devicemanager.New()

	// Read configuration.
	if err := m.readConfig(configFile); err != nil {
		return fmt.Errorf("config file error %s", err)
	}
	// Initiatilze user code.
	m.UserCodes = []*UserCode{
		&UserCode{
			f:    m.timedEvents,
			data: make(chan devicemanager.DeviceData, 10),
		},
		&UserCode{
			f:    m.httpHandler,
			data: make(chan devicemanager.DeviceData, 10),
		},
		&UserCode{
			f:    m.logger,
			data: make(chan devicemanager.DeviceData, 10),
		},
		&UserCode{
			f:    m.modeSleep,
			data: make(chan devicemanager.DeviceData, 10),
		},
	}

	return nil
}

// Run is the main processing loop.
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
			name := m.GetObjectName(got.Object)
			// Set state if object is defined.
			if obj, ok := m.Objects[name]; ok {
				obj.SetState(got.Data)
			}
			// Send event to all user code channels.
			for _, userCode := range m.UserCodes {
				userCode.data <- got
			}

		case err := <-m.DeviceManager.Err:
			log.Printf("device manager error %s", err)
		}
	}
}

// GetObjectName returns the name associated with object address.
func (m *Viki) GetObjectName(address string) string {
	for k, v := range m.Objects {
		if v.Address == address {
			return k
		}
	}
	return ""
}

// SendToObject sends data to the object.
func (m *Viki) ExecObject(name string, data interface{}) error {
	if obj, ok := m.Objects[name]; ok {
		obj.Execute(data)

		// Send state change to all user code channels.
		for _, userCode := range m.UserCodes {
			userCode.data <- devicemanager.DeviceData{
				Data:   data,
				Object: obj.Address,
			}
		}
		return nil
	}
	return fmt.Errorf("unknown object %s", name)
}

// SendToDevice sends data to address on deviceId.
func (m *Viki) SendToDevice(dev string, address string, data interface{}) error {
	if dev, ok := m.DeviceManager.Devices[devicemanager.DeviceId(dev)]; ok {
		dev.Execute(data, address)
		return nil
	}
	return fmt.Errorf("unknown device %s", dev)
}
