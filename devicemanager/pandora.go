package devicemanager

import (
	"fmt"
	"log"
	"net"
	"time"
)

type pandora struct {
	ip   string
	cmd  chan deviceCommand
	quit chan struct{}
	err  chan error
}

func (m *pandora) telnet(ip string, cmd string) error {
	conn, err := net.DialTimeout("tcp", ip, time.Duration(10)*time.Second)
	if err != nil {
		return fmt.Errorf("unable to connect to server %s", ip)
	}
	fmt.Fprintf(conn, cmd)
	return nil
}

func (m *pandora) execute(action string, object string) error {
	log.Printf("pandora: executing %s on %s", action, object)
	cmd := "MSMUSIC\r"
	cmd += "SIPANDORA\r"

	return fmt.Errorf("Ss")
}
func (m *pandora) On() {
	log.Printf("Turn off pandora")
}
func (m *pandora) Off() {
	log.Printf("Turn off pandora")
}

func (m *pandora) Start() {
	log.Printf("starting device pandora...")
	go m.run()
}

func (m *pandora) Shutdown() {
	m.quit <- struct{}{}
}

func (m *pandora) GetErrorChan() <-chan error {
	return m.err
}

func (m *pandora) Execute(action string, object string) {
	m.cmd <- deviceCommand{
		action: action,
		object: object,
	}
}
func (m *pandora) run() {
	for {
		select {
		case cmd := <-m.cmd:
			if err := m.execute(cmd.action, cmd.object); err != nil {
				m.err <- err
			}
		case <-m.quit:
			return
		}
	}
}
