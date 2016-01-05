package main

import (
	"log"
	"viki/tts"
	"time"
)

func main() {
	speaker := tts.New("10.0.0.23:1314", 3)
	speaker.StartSpeaker()

	sp := []string{"hello there", "there is rain", "i am checking it out", "I like ice creams", "golang is concurrency"}
	for _, i := range sp {
		speaker.Speak(i)
		//time.Sleep(1 * time.Second)
	}
	log.Printf("Completed loop")
	time.Sleep(3*time.Second)

	for _, i := range sp {
		speaker.Speak(i)
		//time.Sleep(1 * time.Second)
	}
	for {
	}
}
