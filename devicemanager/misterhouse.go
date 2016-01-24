package devicemanager

import "log"

type mh struct {
	ipPort       string
	host         string
	deviceNumber DeviceNumber
	cmd          chan DeviceCommand
	quit         chan struct{}
	err          chan error
	data         chan DeviceData
}

func (m *mh) execute(action string, object string) error {
	log.Printf("misterhouse %s: executing %s on %s", m.host, action, object)
	return nil
}

func (m *mh) On() {
	log.Printf("Turn off misterhouse")
}
func (m *mh) Off() {
	log.Printf("Turn off misterhouse")
}

func (m *mh) Start() {
	log.Printf("starting device misterhouse...")
	go m.run()
}

func (m *mh) Execute(action string, object string) {
	m.cmd <- DeviceCommand{
		Action: action,
		Object: object,
	}
}
func (m *mh) Shutdown() {
	m.quit <- struct{}{}
}

func (m *mh) GetErrorChan() <-chan error {
	return m.err
}

func (m *mh) run() {
	for {
		select {
		case cmd := <-m.cmd:
			if err := m.execute(cmd.Action, cmd.Object); err != nil {
				m.err <- err
			}
		case <-m.quit:
			return
		}
	}
}
