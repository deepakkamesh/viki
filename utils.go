package viki

import "strings"

func sanitizeState(state string) string {
	recv := strings.Trim(strings.ToLower(state), " ")
	switch recv {
	case "on":
		return "On"
	case "off":
		return "Off"
	default:
		return recv
	}
}
