package main

import (
	"fmt"
	"time"

	"github.com/keep94/sunrise"
)

func main() {
	var s sunrise.Sunrise
	const LATITUDE = float64(37.416969)
	const LONGITUDE = float64(-122.051219)
	s.Around(LATITUDE, LONGITUDE, time.Now())
	s.AddDays(1)
	formatStr := "Jan 2 15:04:05"
	fmt.Printf("Sunrise: %s Sunset: %s\n", s.Sunrise().Format(formatStr), s.Sunset().Format(formatStr))
}
