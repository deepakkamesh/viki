package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.DialTimeout("tcp", "10.0.0.146:23", 2*time.Second)
	fmt.Printf("error %s", err)
	fmt.Print("\nddlklk")
	defer conn.Close()

	//time.Sleep(time.Duration(2) * time.Second)
	cmd := "NS9C\r"
	//cmd += "MSMUSIC\r"
	//cmd += "NS91\r"
	fmt.Fprintf(conn, cmd)

}
