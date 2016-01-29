package main

import (
	"fmt"
	"log"

	"github.com/tarm/serial"
)

func main() {
	fmt.Printf("hello serial")
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 4800}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Printf("Err opening serial")
		return
	}
	b := []byte{0x4, 0x2A, 0x0, 0x6, 0x23, 0x0}
	_, err = s.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
