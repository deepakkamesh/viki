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

// contains returns true if string e exists in slice s.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
