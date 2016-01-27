package main

import (
	"log"
	"os"
	"viki/devicemanager"
)

func main() {
	f, err := os.OpenFile("viki.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Unable to open log file for writing")
	}
	defer f.Close()
	//	log.SetOutput(f)

	log.Print("Logging")
	dm := devicemanager.New()
	dm.StartDeviceManager()
	dm.ExecDeviceCommand(devicemanager.Device_PANDORA, "play", "electronica")
	dm.ExecDeviceCommand(devicemanager.Device_PANDORA, "play", "bollyood")
	dm.ExecDeviceCommand(devicemanager.Device_MISTERHOUSE, "off", "dining")
	dm.ExecDeviceCommand(devicemanager.Device_MISTERHOUSE, "on", "living")
	for {
	}
}
