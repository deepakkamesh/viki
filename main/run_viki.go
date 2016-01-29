package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"viki"
)

func main() {

	v := viki.New()
	fmt.Println("Starting Viki version:", v.Version)

	// Setup flags.
	configFile := flag.String("config_file", "../objects.conf", "Config file for objects")
	logFile := flag.String("log_file", "viki.log", "log file path")
	flag.String("festival_ipport", "10.0.0.23:1314", "Ip:Port of festival server")
	flag.String("http_listen_port", "2233", "Port number of the http server")
	flag.String("x10_tty", "/dev/ttyUSB0", "tty device for x10 controller")
	flag.Parse()

	// Setup log file.
	f, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Unable to open log %s file for writing %s", *logFile, err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Init and run viki
	if err := v.Init(*configFile); err != nil {
		log.Fatalf("Fatal Error: %s\n", err)
	}
	v.Run()
}
