package devicemanager

import "log"

type HttpHandler struct {
	deviceNumber DeviceNumber
	cmd          chan DeviceCommand
	quit         chan struct{}
	err          chan error
	data         chan DeviceData
}

func getRealIP() (string, error) {
	return "10.0.0.1", nil
}

func (m *HttpHandler) execute(action string, object string) (string, error) {
	log.Printf("HttpHandler: executing %s on %s", action, object)
	return getRealIP()
}

func (m *HttpHandler) On() {
	log.Printf("Turn off HttpHandler")
}
func (m *HttpHandler) Off() {
	log.Printf("Turn off HttpHandler")
}

func (m *HttpHandler) Start() {
	log.Printf("starting device HttpHandler...")
	go m.run()
}

func (m *HttpHandler) Shutdown() {
	m.quit <- struct{}{}
}

func (m *HttpHandler) GetErrorChan() <-chan error {
	return m.err
}

func (m *HttpHandler) Execute(action string, object string) {
	m.cmd <- DeviceCommand{
		Action: action,
		Object: object,
	}
}
func (m *HttpHandler) run() {
	for {
		select {
		case cmd := <-m.cmd:
			data, err := m.execute(cmd.Action, cmd.Object)
			if err != nil {
				m.err <- err
				continue
			}
			m.data <- DeviceData{
				DeviceNumber: m.deviceNumber,
				Data:         data,
				Object:       cmd.Object,
			}
		case <-m.quit:
			return
		}
	}
}
