/* Package viki is a intelligent extensible home automation framework.
 */
package viki

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/deepakkamesh/viki/devicemanager/device"
	"github.com/deepakkamesh/viki/objectmanager"
	"github.com/golang/glog"
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

	// Initialization order is important since readConfig relies
	// on an initialized ObjectManager and DeviceManager.
	m.ObjectManager = objectmanager.New()
	m.DeviceManager = devicemanager.New(m.ObjectManager)

	if err := m.readConfig(configFile); err != nil {
		return fmt.Errorf("failed to read config file %v", err)
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

	// Process messages from devices.
	for {
		select {
		case got := <-m.DeviceManager.Data:
			_, o := m.ObjectManager.GetObjectByAddress(got.Address)
			if o == nil {
				glog.Warningf("Object with address %s does not exit ", got.Address)
				continue
			}
			o.SetState(got.Data)
			// Send event to all user code channels.
			for _, userChan := range m.userChannels {
				userChan.data <- got
			}
		case err := <-m.DeviceManager.Err:
			glog.Errorf("Device manager error %s", err)
		}
	}
}

// Viki Helper Functions for usercode.

// Do sends data to the object with name.
func (m *Viki) Do(name string, data interface{}) error {

	a, o := m.ObjectManager.GetObjectByName(name)
	if o == nil {
		glog.Errorf("object with name %v not found", name)
		return fmt.Errorf("object with name %v not found", name)
	}

	if err := o.Execute(data); err != nil {
		glog.Errorf("failed to execute %v", err)
		return err
	}

	// Send state change to all user code channels.
	for _, userChan := range m.userChannels {
		userChan.data <- devicemanager.DeviceData{
			Data:    data,
			Address: a,
		}
	}
	return nil
}
