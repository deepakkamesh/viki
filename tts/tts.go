// Package tts does buffered text to speech conversion using misterhouse.
package tts

import (
	"fmt"
	"log"
	"net"
	"time"
)

type tts struct {
	ipPort  string
	timeout int64
	speaker chan string
	quit    chan struct{}
	conn    net.Conn
}

func (m *tts) Speak(text string) {
	m.speaker <- text
}

func (m *tts) mockSpeak(text string) error {
	log.Printf("recieved %s", text)
	return nil
}

/*
// speak sends a voice command to MisterHouse.
func (m *tts) speak(text string) error {
	conn, err := net.DialTimeout("tcp", m.ipPort, time.Duration(m.timeout)*time.Second)
	if err != nil {
		return fmt.Errorf("Unable to connect to %s:%s", m.ipPort, err)
	}
	defer conn.Close()
	fmt.Fprintf(conn, "(tts_text \"%s\" nil)(quit)", text)
	return nil
}

*/
// speak sends a voice command to MisterHouse.
func (m *tts) speak(text string) error {
	var errW, err error
	r := 3
	// If conn is not initialized, attempt to connect.
	if m.conn == nil {
		log.Printf("tts.go: connection not initialized. Starting tts connection...")
		if m.conn, err = net.DialTimeout("tcp", m.ipPort, time.Duration(m.timeout)*time.Second); err != nil {
			return fmt.Errorf("tts.go: unable to connect to %s:%s", m.ipPort, err)
		}
	}

	// Try to write to socket, if not try to connect and write again.
	for _, errW = fmt.Fprintf(m.conn, "(tts_text \"%s\" nil)", text); errW != nil && r > 0; r -= 1 {
		log.Printf("tts.go: tts socket write failed. Retrying connection...")
		if m.conn, err = net.DialTimeout("tcp", m.ipPort, time.Duration(m.timeout)*time.Second); err != nil {
			return fmt.Errorf("tts.go: unable to connect to %s:%s", m.ipPort, err)
		}
	}
	if errW != nil {
		return fmt.Errorf("tts.go: unable to write to socket %s:%s", m.ipPort, err)
	}
	return nil
}

// Shutdown stops the speech loop.
func (m *tts) ShutDown() {
	m.conn.Close()
	m.quit <- struct{}{}
}

// StartSpeaker starts the speech loop.
func (m *tts) StartSpeaker() {
	go m.runSpeech()
}

// runSpeech runs the speech loop.
func (m *tts) runSpeech() {
	log.Printf("tts.go: starting speech loop...")
	for {
		select {
		case tosay := <-m.speaker:
			if err := m.speak(tosay); err != nil {
				log.Printf("tts.go: %s", err)
			}
			log.Printf("tts.go: speaking %s", tosay)
		case <-m.quit:
			return
		}
	}
}

// New returns an initialized TTS module.
func New(ipPort string, timeout int64) *tts {
	return &tts{
		ipPort:  ipPort,
		timeout: timeout,
		speaker: make(chan string, 10),
		quit:    make(chan struct{}),
		conn:    nil,
	}
}
