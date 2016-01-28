package devicemanager

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HttpHandler struct {
	deviceNumber DeviceNumber
	cmd          chan DeviceData
	quit         chan struct{}
	err          chan error
	data         chan DeviceData
}

func (m *HttpHandler) On() {
}
func (m *HttpHandler) Off() {
}

func (m *HttpHandler) Start() error {
	log.Printf("starting device HttpHandler...")
	http.HandleFunc("/object/", m.handleObject)
	fl := flag.Lookup("http_listen_port")
	port := fl.Value.String()
	go http.ListenAndServe(":"+port, nil)
	go m.run()
	return nil
}

func (m *HttpHandler) Shutdown() {
	m.quit <- struct{}{}
}

func (m *HttpHandler) Execute(action interface{}, object string) {
	m.cmd <- DeviceData{
		Data:   action,
		Object: object,
	}
}

func (m *HttpHandler) run() {
	for {
		select {
		case <-m.cmd:
			continue
		case <-m.quit:
			return
		}
	}
}

func (m *HttpHandler) handleObject(w http.ResponseWriter, r *http.Request) {
	req := strings.Split(r.URL.Path[1:], "/")
	if req[0] == "object" {
		m.data <- DeviceData{
			DeviceNumber: m.deviceNumber,
			Data:         req[1:],
			Object:       "http",
		}
		fmt.Fprintf(w, "Setting %s on %s", req[2], req[1])
		log.Printf("recieved http request %s %s", req[2], req[1])
	}
}
