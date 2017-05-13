package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/deepakkamesh/viki"
	"github.com/golang/glog"
)

var (
	buildtime string
	githash   string
)

func main() {

	// Setup flags.
	configFile := flag.String("config_file", "objects.conf", "Config file for objects")
	version := flag.Bool("version", false, "display version")

	flag.String("festival_ipport", "127.0.0.1:1314", "Ip:Port of festival server")
	flag.String("mochad_ipport", "127.0.0.1:1099", "Ip:Port of mochad server")
	flag.String("graphite_ipport", "", "Ip:port of graphite server.")
	flag.String("http_listen_port", "2233", "Port number of the http server")
	flag.String("x10_tty", "/dev/ttyUSB0", "tty device for x10 controller")
	flag.String("resource", "./resources", "path to the resources folder")
	flag.Float64("lat", 0, "latitude coordinate")
	flag.Float64("long", 0, "longitude coordinate")
	flag.String("mg_domain", "", "Domain name for mailgun")
	flag.String("mg_apikey", "", "Api key for mailgun")
	flag.String("mg_pubkey", "", "Public API key for mailgun")
	flag.String("email_alert_list", "", "comma separated list of people to alert by email for events")
	flag.Bool("ssl", false, "listen only on https")
	flag.Parse()

	// Print version and exit.
	if *version {
		fmt.Printf("Version commit hash %s\n", githash)
		fmt.Printf("Build date %s\n", buildtime)
		os.Exit(0)
	}

	// Init and run viki
	v := viki.New(githash)
	glog.Infof("Starting Viki version:%v", v.Version)
	if err := v.Init(*configFile); err != nil {
		glog.Fatalf("Fatal Error: %s\n", err)
	}
	v.Run()
}
