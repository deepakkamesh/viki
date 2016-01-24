package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)
import "time"

import "net"

func main() {
	conn, err := net.DialTimeout("tcp", "ifconfig.me:80", time.Duration(10)*time.Second)
	if err != nil {
		fmt.Printf("unable to connect to server ifconfig.me %s", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	resp, err := http.Get("http://ifconfig.me/ua")
	fmt.Printf("%+v", resp)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s", body)
}
