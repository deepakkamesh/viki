package viki

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/deepakkamesh/viki/devicemanager"
	"github.com/golang/glog"
)

const Graphite_PREFIX string = "viki"

func (m *Viki) MyLogger(c chan devicemanager.DeviceData) {

	glog.Infof("Starting user routine MyLogger...")
	defer glog.Infof("Shutting down user routine MyLogger...")

	graphiteIpPort := flag.Lookup("graphite_ipport").Value.String()

	for {
		// Wait to recieve any events.
		got := <-c

		name, object := m.ObjectManager.GetObjectByAddress(got.Address)

		// TODO: this might change if type is different.
		state, ok := got.Data.(string)
		if !ok {
			continue
		}

		glog.V(1).Infof("Event: %s(%s),%s", name, got.DeviceId, state)

		// Log to graphite server if enabled.
		if graphiteIpPort == "" {
			continue
		}
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
			glog.Warningf("unable to dial graphite %s", err)
			continue
		}
		if _, err := fmt.Fprintf(conn, "%s\n", metric); err != nil {
			glog.Warningf("unable to send metric to graphite %s", err)
		}
		continue
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
