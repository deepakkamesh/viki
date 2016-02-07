package viki

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/deepakkamesh/viki/devicemanager"
)

func (m *Viki) httpHandler(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine httphandler...")

	// Build object list for nlp.
	objs := []*nlpMatch{}
	for k, _ := range m.Objects {
		objs = append(objs, &nlpMatch{
			object: k,
			weight: 0,
		})
	}

	// Read the templates.
	var res string
	fl := flag.Lookup("resource")
	res = fl.Value.String()

	tplIdx, err := template.ParseFiles(res + "/index.html")
	if err != nil {
		log.Fatalf("error reading file %s", err)
	}
	tick := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-tick.C:
			// Build and send objects and state to http device.
			idxPage, err := buildIndexPage(m.Objects, tplIdx)
			if err != nil {
				log.Printf("build index page failed %s", err)
				continue
			}
			m.SendToDevice("httphandler", "idxpage", idxPage)

		// Channel to recieve any events.
		case got := <-c:
			switch got.Object {
			case "http_cmd":
				d, _ := got.Data.([]string)
				state := sanitizeState(d[1])
				if err := m.ExecObject(d[0], state); err != nil {
					log.Printf("recieved unknown object %s", d[0])
					continue
				}
				m.ExecObject("speaker", "Executing command")

			case "http_qry":
				d, _ := got.Data.(string)
				res := matchObject(objs, d)
				if act, err := matchAction(d); err == nil {
					for _, i := range res {
						m.ExecObject(i.object, act)
					}
				}
			}
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
	return o.Objects[i].Object < o.Objects[j].Object
}

func (o idxtpl) Swap(i, j int) {
	o.Objects[i], o.Objects[j] = o.Objects[j], o.Objects[i]
}

// buildIndexPage constructs the index page for the web interface.
func buildIndexPage(o map[string]*Object, tplIdx *template.Template) (string, error) {

	objs := idxtpl{}

	buf := new(bytes.Buffer)

	// Add object to list if its not hidden.
	for k, v := range o {
		if !v.CheckTag("web_hidden") {
			state := "NA"
			if v.State != nil {
				state = v.State.(string)
			}
			objs.Objects = append(objs.Objects, objtpl{
				Object: k,
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

/* Natural Language Processing */

type nlpMatch struct {
	object string
	weight int
}

// matchAction matches the state/action requested in ivr.
func matchAction(ivr string) (string, error) {
	input_words := strings.Split(strings.ToLower(ivr), " ")
	actions := map[string]string{
		"on":  "On",
		"off": "Off",
	}

	for key, value := range actions {
		if contains(input_words, key) {
			return value, nil
		}
	}
	return "", fmt.Errorf("no matching action")
}

// Functions to handle Sort.
type byWeight []*nlpMatch

func (a byWeight) Len() int {
	return len(a)
}
func (a byWeight) Less(i, j int) bool {
	return a[i].weight < a[j].weight
}
func (a byWeight) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// matchObject matches sentence with objs.
func matchObject(objs []*nlpMatch, sentence string) []*nlpMatch {

	input_words := strings.Split(sentence, " ")

	object_weights := make(map[string]int)

	total_words := 0
	for i := 0; i < len(objs); i++ {
		words := strings.Split(objs[i].object, " ")
		for _, word := range words {
			object_weights[word]++
			total_words++
		}
	}
	// Refactor final object_weights - higher the weight, rarer the word
	for key, _ := range object_weights {
		object_weights[key] = total_words - object_weights[key]
	}

	//Calculate the relative hit rate of objects
	for i := 0; i < len(objs); i++ {
		words := strings.Split(objs[i].object, " ")
		weight := 0
		for _, word := range words {
			if contains(input_words, word) {
				weight += object_weights[word]
			}
		}
		objs[i].weight = weight
	}

	// Sort by weights in reverse.
	sort.Sort(sort.Reverse(byWeight(objs)))

	// Return only the highest matching entries.
	matcher := objs[0].weight
	for i := 1; i < len(objs); i++ {
		if objs[i].weight != matcher {
			return objs[:i]
		}
	}
	return nil
}
