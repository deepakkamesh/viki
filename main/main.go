package main

import (
	"fmt"
	"viki"
)

func main() {

	v := viki.New()
	fmt.Println(v.Version)
	v.Init()
	v.Run()

	for{}
}
