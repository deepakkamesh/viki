package viki

import (
	"fmt"
	"sort"
	"strings"

	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/golang/glog"
)

func (m *Viki) MyHTTPHandler(c chan devicemanager.DeviceData) {

	glog.Infof("Starting user routine MyHTTPHandler...")
	defer glog.Infof("Shutting down routine MyHTTPHandler")

	// Build object list for nlp.
	objs := []*nlpMatch{}
	for _, o := range m.ObjectManager.Objects {
		objs = append(objs, &nlpMatch{
			object: o.Name,
			weight: 0,
		})
	}

	for {
		select {
		case got := <-c:
			switch got.Address {
			case "http_cmd":
				d, _ := got.Data.([]string)
				state := sanitizeState(d[1])
				if err := m.Do(d[0], state); err != nil {
					glog.Errorf("Error executing on %s %v", d[0], err)
					continue
				}

			case "http_qry":
				d, _ := got.Data.(string)
				res := matchObject(objs, d)
				if act, err := matchAction(d); err == nil {
					for _, i := range res {
						if err := m.Do(i.object, act); err != nil {
							glog.Errorf("Error executing query %s on %s: %v", d, i.object, err)
						}
					}
				}
			}
		}
	}
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
