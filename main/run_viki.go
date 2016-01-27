package main

import (
	"fmt"
	"os"
	"viki"
)

func main() {

	v := viki.New()
	fmt.Println("main run", v.Version)
	if err := v.Init(); err != nil {
		fmt.Printf("Fatal Error: %s\n", err)
		os.Exit(1)
	}
	v.Run()
	for {
	}
}
