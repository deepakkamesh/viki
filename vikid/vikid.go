package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/deepakkamesh/viki"
)

var (
	buildtime string
	githash   string
)

func main() {

	// Setup flags.
	configFile := flag.String("config_file", "objects.conf", "Config file for objects")
	logFile := flag.String("log_file", "viki.log", "log file path")
	logStdOut := flag.Bool("log_stdout", true, "log to std out only")
	version := flag.Bool("version", false, "display version")
	flag.String("festival_ipport", "10.0.0.102:1314", "Ip:Port of festival server")
	flag.String("mochad_ipport", "10.0.0.102:1099", "Ip:Port of mochad server")
	flag.String("graphite_ipport", "", "Ip:port of graphite server.")
	flag.String("http_listen_port", "2233", "Port number of the http server")
	flag.String("x10_tty", "/dev/ttyUSB0", "tty device for x10 controller")
	flag.String("resource", "./resources", "path to the resources folder")
	flag.String("log", "./logs", "path to the logs folder")
	flag.Float64("lat", 37.416969, "latitude coordinate")
	flag.Float64("long", -122.051219, "longitude coordinate")
	flag.Bool("ssl", false, "listen only on https")
	flag.Parse()

	// Print version and exit.
	if *version {
		fmt.Printf("Version commit hash %s\n", githash)
		fmt.Printf("Build date %s\n", buildtime)
		os.Exit(0)
	}

	if !*logStdOut {
		// Setup log file.
		f, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Fatalf("Unable to open log %s file for writing %s", *logFile, err)
		}
		log.SetOutput(f)
		defer f.Close()
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Init and run viki
	v := viki.New(githash)
	fmt.Println("Starting Viki version:", v.Version)
	if err := v.Init(*configFile); err != nil {
		log.Fatalf("Fatal Error: %s\n", err)
	}
	v.Run()
}
