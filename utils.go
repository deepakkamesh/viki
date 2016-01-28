package viki

import "strings"

func sanitizeState(state string) string {

	switch strings.Trim(strings.ToLower(state), " ") {
	case "on":
		return "On"
	case "off":
		return "Off"
	}
	return ""
}
