package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

func main() {
	type objTpl struct {
		Object string
		State  string
		some   int
	}
	st := []objTpl{
		objTpl{
			Object: "living",
			State:  "On",
		},
		objTpl{
			Object: "dinin",
			State:  "Off",
		},
	}
	tpl, err := template.ParseFiles("/Users/dkg/Projects/golang/src/github.com/deepakkamesh/viki/resources/object.html")
	if err != nil {
		log.Printf("Template error %s", err)
	}
	buf := new(bytes.Buffer)
	for _, i := range st {
		_ = tpl.Execute(buf, i)
	}

	fmt.Println(buf.String())
}
