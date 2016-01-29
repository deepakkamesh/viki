package devicemanager

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Unique Device Number.
const Device_HTTPHANDLER DeviceId = "httphandler"

type HttpHandler struct {
	deviceId DeviceId
	in       chan DeviceData
	quit     chan struct{}
	err      chan error
	out      chan DeviceData
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
	m.in <- DeviceData{
		Data:   action,
		Object: object,
	}
}

func (m *HttpHandler) run() {
	for {
		select {
		case <-m.in:
			continue
		case <-m.quit:
			return
		}
	}
}

func (m *HttpHandler) handleObject(w http.ResponseWriter, r *http.Request) {
	req := strings.Split(r.URL.Path[1:], "/")
	if req[0] == "object" {
		m.out <- DeviceData{
			DeviceId: m.deviceId,
			Data:     req[1:],
			Object:   "http",
		}
		fmt.Fprintf(w, "Setting %s on %s", req[2], req[1])
		log.Printf("recieved http request %s %s", req[2], req[1])
	}
}
