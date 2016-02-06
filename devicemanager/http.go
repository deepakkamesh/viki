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
	idxPage  string
}

func (m *HttpHandler) On() {
}
func (m *HttpHandler) Off() {
}

func (m *HttpHandler) Start() error {
	log.Printf("starting device HttpHandler...")
	http.HandleFunc("/object/", m.handleObject)
	http.HandleFunc("/q/", m.handleQuery)
	http.HandleFunc("/", m.handleIndex)
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
		case got := <-m.in:
			if got.Object == "idxpage" {
				m.idxPage = got.Data.(string)
			}
		case <-m.quit:
			return
		}
	}
}

func (m *HttpHandler) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%s", m.idxPage)
}

func (m *HttpHandler) handleObject(w http.ResponseWriter, r *http.Request) {
	req := strings.Split(r.URL.Path[1:], "/")
	if len(req) < 3 {
		fmt.Fprintf(w, "Error: Use format object/<name>/<cmd>")
	}
	m.out <- DeviceData{
		DeviceId: m.deviceId,
		Data:     req[1:],
		Object:   "http_cmd",
	}
	log.Printf("recieved http request %s %s", req[2], req[1])
	http.Redirect(w, r, "/", 302)
}

func (m *HttpHandler) handleQuery(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Path[3:]

	m.out <- DeviceData{
		DeviceId: m.deviceId,
		Data:     q,
		Object:   "http_qry",
	}
	fmt.Fprintf(w, "Executing query %s", q)
	log.Printf("recieved http request %s", q)
}
