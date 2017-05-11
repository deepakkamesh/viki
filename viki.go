/* Package viki is a intelligent extensible home automation framework.
 */
package viki

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/deepakkamesh/viki/devicemanager/device"
	"github.com/deepakkamesh/viki/objectmanager"
)

type userChannel struct {
	fName string                        // Name of method.
	data  chan devicemanager.DeviceData // Allocated channel.
}

type Viki struct {
	Version       string
	ObjectManager *objectmanager.ObjectManager
	DeviceManager *devicemanager.DeviceSettings
	userChannels  []userChannel
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
			dev device.Device
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
		m.ObjectManager.AddObject(c[1], c[0], dev, tags)

	}
	return nil
}

func New(ver string) *Viki {
	return &Viki{
		Version: ver,
	}
}

func (m *Viki) Init(configFile string) error {

	// Initialize device manager and object manager.
	m.ObjectManager = objectmanager.New()
	m.DeviceManager = devicemanager.New(m.ObjectManager)

	// Read configuration.
	if err := m.readConfig(configFile); err != nil {
		return fmt.Errorf("config file error %s", err)
	}

	return nil
}

// Run is the main processing loop.
func (m *Viki) Run() {
	// Start Device Manager.
	m.DeviceManager.StartDeviceManager()

	// Start user code.
	// Uses reflection to enumerate methods starting with My*.
	typ := reflect.TypeOf(m)
	for i := 0; i < typ.NumMethod(); i++ {
		if regexp.MustCompile("^My(.+)").MatchString(typ.Method(i).Name) {
			recv := make(chan devicemanager.DeviceData, 10)
			m.userChannels = append(m.userChannels, userChannel{
				fName: typ.Method(i).Name,
				data:  recv,
			})
			// Run the user code in a goroutine.
			go reflect.ValueOf(m).Method(i).Call([]reflect.Value{
				reflect.ValueOf(recv),
			})
		}
	}

	// Run the main processing loop.
	for {
		select {
		case got := <-m.DeviceManager.Data:
			_, o := m.ObjectManager.GetObjectByAddress(got.Object)
			if o != nil {
				// Set state if object is defined.
				o.SetState(got.Data)
				// Send event to all user code channels.
				for _, userChan := range m.userChannels {
					userChan.data <- got
				}
				continue
			}
			log.Printf("object %s does not exit ", got.Object)

		case err := <-m.DeviceManager.Err:
			log.Printf("device manager error %s", err)
		}
	}
}

// execObject sends data to the object.
func (m *Viki) execObject(name string, data interface{}) error {
	if err := m.ObjectManager.Exec(name, data); err != nil {
		return err
	}

	a, _ := m.ObjectManager.GetObjectByName(name)

	// Send state change to all user code channels.
	for _, userChan := range m.userChannels {
		userChan.data <- devicemanager.DeviceData{
			Data:   data,
			Object: a,
		}
	}
	return nil
}

// SendToDevice sends data to address on deviceId.
func (m *Viki) sendToDevice(dev string, address string, data interface{}) error {
	if dev, ok := m.DeviceManager.Devices[devicemanager.DeviceId(dev)]; ok {
		dev.Execute(data, address)
		return nil
	}
	return fmt.Errorf("unknown device %s", dev)
}
