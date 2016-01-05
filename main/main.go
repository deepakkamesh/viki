package main

import (
	"flag"
	"log"
	"viki/tts"
)

func main() {
	// Grab all the flags.
	//logFile := flag.String("log", "/var/log/viki.log", "The location of the log file for viki")
	ttsTest := flag.String("tts_test", "Hello, this is viki", "Test words to speak")
	//mhIP := flag.String("mh_ip", "10.0.0.23", "IP address of MisterHouse Server")
	ttsIpPort := flag.String("mh_tts_port", "10.0.0.23:1314", "Port number of the TTS server of misterhouse")
	flag.Parse()
	// Initialize the backends.
	speaker := tts.New(*ttsIpPort, 10)

	if err := speaker.Speak(*ttsTest); err != nil {
		log.Printf("Error %s", err)
	}
}
