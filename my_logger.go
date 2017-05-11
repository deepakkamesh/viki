package viki

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/deepakkamesh/viki/devicemanager"
)

const Graphite_PREFIX string = "viki"

func (m *Viki) MyLogger(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine logger...")

	graphiteIpPort := flag.Lookup("graphite_ipport").Value.String()

	for {
		// Wait to recieve any events.
		got := <-c

		name, object := m.ObjectManager.GetObjectByAddress(got.Object)
		// TODO: this might change if type is different.
		state, ok := got.Data.(string)
		if !ok {
			continue
		}

		if object != nil {

			log.Printf("event: %s(%s),%s", name, got.DeviceId, state)

			// Log to graphite server if enabled.
			if graphiteIpPort != "" {
				metricPrefix := Graphite_PREFIX
				if object.CheckTag("motion") {
					metricPrefix += ".motion"
				} else if object.CheckTag("door") {
					metricPrefix += ".door"
				} else if object.CheckTag("appliance") {
					metricPrefix += ".appliance"
				} else {
					continue
				}
				metric := formatMetric(metricPrefix, name, state)
				conn, err := net.DialTimeout("tcp", graphiteIpPort, time.Duration(2)*time.Second)
				if err != nil {
					log.Printf("unable to dial graphite %s", err)
					continue
				}
				if _, err := fmt.Fprintf(conn, "%s\n", metric); err != nil {
					log.Printf("unable to send metric to graphite %s", err)
				}
			}
			continue
		}
		log.Printf("Got data from unknown object %s %s\n", got.Object, state)
	}
}

func formatMetric(prefix, name, state string) string {

	metric := prefix
	metric += "." + strings.Replace(name, " ", "_", -1)
	metric += " " + translateState(state)
	metric += " " + strconv.FormatInt(time.Now().Unix(), 10)

	return metric
}

func translateState(st string) string {
	switch st {
	case "On":
		return "1"
	case "Off":
		return "0"
	case "Open":
		return "1"
	case "Closed":
		return "0"
	default:
		return "0"
	}
}
