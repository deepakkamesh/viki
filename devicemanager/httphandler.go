package devicemanager

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Fullfillment Request.
type FulfillmentRequest struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Result    struct {
		Source           string `json:"source"`
		ResolvedQuery    string `json:"resolvedQuery"`
		Action           string `json:"action"`
		ActionIncomplete bool   `json:"actionIncomplete"`
		Parameters       struct {
			Object string `json:"object"`
			State  string `json:"state`
		} `json:"parameters"`
		Contexts []struct {
			Name       string `json:"name"`
			Parameters struct {
				Name string `json:"name"`
			} `json:"parameters"`
			Lifespan int `json:"lifespan"`
		} `json:"contexts"`
		Metadata struct {
			IntentID   string `json:"intentId"`
			IntentName string `json:"intentName"`
		} `json:"metadata"`
		Fulfillment struct {
			Speech string `json:"speech"`
		} `json:"fulfillment"`
	} `json:"result"`
	Status struct {
		Code      int    `json:"code"`
		ErrorType string `json:"errorType"`
	} `json:"status"`
}

// Unique Device Number.
const Device_HTTPHANDLER DeviceId = "httphandler"

type httphandler struct {
	in      chan DeviceData
	quit    chan struct{}
	err     chan error
	out     chan DeviceData
	idxPage string
}

// NewDevHttpHandler returns a new initialized http handler.
func (m *DeviceSettings) NewDeviceHttpHandler(out chan DeviceData, err chan error) (DeviceId, Device) {
	return Device_HTTPHANDLER, &httphandler{
		in:   make(chan DeviceData, 10),
		quit: make(chan struct{}),
		err:  err,
		out:  out,
	}
}

func (m *httphandler) On() {
}
func (m *httphandler) Off() {
}

func (m *httphandler) Start() error {
	log.Printf("starting device HttpHandler...")
	http.HandleFunc("/object/", m.handleObject)         // Handler for commands on objects.
	http.HandleFunc("/q/", m.handleQuery)               // Handler for  queries (nlp).
	http.HandleFunc("/googlehome/", m.handleGoogleHome) // Handler for google home.
	http.HandleFunc("/", m.handleIndex)

	port := flag.Lookup("http_listen_port").Value.String()
	res := flag.Lookup("resource").Value.String()
	ssl := flag.Lookup("ssl").Value.String()

	if ssl == "true" {
		go http.ListenAndServeTLS(":"+port, res+"/server.crt", res+"/server.key", nil)
	} else {
		go http.ListenAndServe(":"+port, nil)
	}
	go m.run()
	return nil
}

func (m *httphandler) Shutdown() {
	m.quit <- struct{}{}
}

func (m *httphandler) Execute(action interface{}, object string) {
	m.in <- DeviceData{
		Data:   action,
		Object: object,
	}
}

func (m *httphandler) run() {
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

// handleGoogleHome is the http handler for Google Home integration.
func (m *httphandler) handleGoogleHome(w http.ResponseWriter, r *http.Request) {
	var msg FulfillmentRequest
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("Failed to decode json fulfllment request")
	}

}

// handleIndex is the http handler for the index page.
func (m *httphandler) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%s", m.idxPage)
}

// handleObject is the http handler for the object command.
func (m *httphandler) handleObject(w http.ResponseWriter, r *http.Request) {
	req := strings.Split(r.URL.Path[1:], "/")
	if len(req) < 3 {
		fmt.Fprintf(w, "Error: Use format object/<name>/<cmd>")
	}
	m.out <- DeviceData{
		DeviceId: Device_HTTPHANDLER,
		Data:     req[1:],
		Object:   "http_cmd",
	}
	log.Printf("recieved http request %s %s", req[2], req[1])
	http.Redirect(w, r, "/", 302)
}

// handleQuery is the http handler for natural language.
func (m *httphandler) handleQuery(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(r.URL.Path[3:])

	m.out <- DeviceData{
		DeviceId: Device_HTTPHANDLER,
		Data:     q,
		Object:   "http_qry",
	}
	fmt.Fprintf(w, "executing  %s", q)
	log.Printf("recieved http request %s", q)
}
