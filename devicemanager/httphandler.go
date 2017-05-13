package devicemanager

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/deepakkamesh/viki/devicemanager/device"
	"github.com/deepakkamesh/viki/objectmanager"
	"github.com/golang/glog"
)

// Google Home Fullfillment Response.
type PermissionsRequest struct {
	OptContext  string   `json:"opt_context"`
	Permissions []string `json:"permissions"`
}
type Google struct {
	ExpectUserResponse bool `json:"expect_user_response"`
	IsSsml             bool `json:"is_ssml"`
	//	PermissionsRequest PermissionsRequest `json:"permissions_request"` // Only needed if there are perms.
}
type Data struct {
	Google Google `json:"google"`
}
type FulfillmentResponse struct {
	Speech      string   `json:"speech"`
	DisplayText string   `json:"displayText"`
	Data        Data     `json:"data"`
	ContextOut  []string `json:"contextOut"`
	Source      string   `json:"source"`
}

// Google Home Fullfillment Request.
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
	in   chan DeviceData
	quit chan struct{}
	err  chan error
	out  chan DeviceData
	om   *objectmanager.ObjectManager
}

// NewDevHttpHandler returns a new initialized http handler.
func (m *DeviceSettings) NewDeviceHttpHandler(out chan DeviceData, err chan error, om *objectmanager.ObjectManager) (DeviceId, device.Device) {
	return Device_HTTPHANDLER, &httphandler{
		in:   make(chan DeviceData, 10),
		quit: make(chan struct{}),
		err:  err,
		out:  out,
		om:   om,
	}
}

func (m *httphandler) On() {
}
func (m *httphandler) Off() {
}

func (m *httphandler) Start() error {
	glog.Infof("starting device HttpHandler...")
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
		Data:    action,
		Address: object,
	}
}

func (m *httphandler) run() {
	for {
		select {
		case <-m.in:
		case <-m.quit:
			return
		}
	}
}

/* Building the index page */
type objtpl struct {
	Object string
	State  string
	Ro     bool
}
type idxtpl struct {
	Objects []objtpl
}

func (o idxtpl) Len() int {
	return len(o.Objects)
}

func (o idxtpl) Less(i, j int) bool {

	if o.Objects[i].Ro == o.Objects[j].Ro {
		return o.Objects[i].Object < o.Objects[j].Object
	}
	// Sort by state Ro or Not.
	return !o.Objects[i].Ro && o.Objects[j].Ro
}

func (o idxtpl) Swap(i, j int) {
	o.Objects[i], o.Objects[j] = o.Objects[j], o.Objects[i]
}

// buildIndexPage constructs the index page for the web interface.
func buildIndexPage(o []*objectmanager.Object, tplIdx *template.Template) (string, error) {

	objs := idxtpl{}

	buf := new(bytes.Buffer)

	// Add object to list if its not hidden.
	for _, v := range o {
		if !v.CheckTag("web_hidden") {
			state := "NA"
			if v.State != nil {
				state = v.State.(string)
			}
			objs.Objects = append(objs.Objects, objtpl{
				Object: v.Name,
				State:  state,
				Ro:     v.CheckTag("web_ro"),
			})
		}
	}
	sort.Sort(objs)
	if err := tplIdx.Execute(buf, objs); err != nil {
		return "", fmt.Errorf("error parsing template %s", err)
	}

	return buf.String(), nil
}

// handleIndex is the http handler for the index page.
func (m *httphandler) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Read the templates.
	var res string
	fl := flag.Lookup("resource")
	res = fl.Value.String()

	tplIdx, err := template.ParseFiles(res + "/index.html")
	if err != nil {
		glog.Fatalf("error reading file %s", err)
	}
	// Build and send objects and state to http device.
	idxPage, err := buildIndexPage(m.om.Objects, tplIdx)
	if err != nil {
		glog.Infof("build index page failed %s", err)
	}

	fmt.Fprintf(w, "%s", idxPage)
}

// handleGoogleHome is the http handler for Google Home integration.
// Google Home integration works via api.ai. Any new objects/states should
// be configured within api.ai as well.
func (m *httphandler) handleGoogleHome(w http.ResponseWriter, r *http.Request) {
	var msg FulfillmentRequest
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		glog.Infof("Failed to decode json fulfllment request")
	}

	glog.Infof("Got request from google home %v", msg.Result.Parameters)
	object := msg.Result.Parameters.Object
	state := msg.Result.Parameters.State

	m.out <- DeviceData{
		DeviceId: Device_HTTPHANDLER,
		Data:     []string{object, state},
		Address:  "http_cmd",
	}

	// Build a response back.
	resp := FulfillmentResponse{
		Speech:      "Ok,I have turned " + state + " " + object,
		DisplayText: "Ok,All Done",
		Data: Data{
			Google{
				ExpectUserResponse: false,
				IsSsml:             false,
			},
		},
		ContextOut: []string{},
		Source:     "Viki",
	}

	b, err := json.Marshal(resp)
	if err != nil {
		glog.Fatalf("Failed to marshal response: %v", err)
	}

	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Google-Assistant-API-Version", "v1")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
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
		Address:  "http_cmd",
	}
	glog.Infof("recieved http request %s %s", req[2], req[1])
	http.Redirect(w, r, "/", 302)
}

// handleQuery is the http handler for natural language.
func (m *httphandler) handleQuery(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(r.URL.Path[3:])

	m.out <- DeviceData{
		DeviceId: Device_HTTPHANDLER,
		Data:     q,
		Address:  "http_qry",
	}
	fmt.Fprintf(w, "executing  %s", q)
	glog.Infof("recieved http request %s", q)
}
