package main

import "fmt"

type obj struct {
	a int
	b interface{}
}

type state struct {
	b string
}

func (m *state) GetState() string {
	return m.b
}

func (m *state) SetState(st string) {
	m.b = st
}

func main() {
	a := obj{
		a: 3,
		b: &state{},
	}
	c, _ := a.b.(*state)
	c.SetState("dddddd")
	fmt.Println(c.GetState())
}
