package viki

// getMochadState returns the string value of state of the object named name
// after asserting.
func (m *Viki) getMochadState(name string) string {
	if v, ok := m.Objects[name]; ok {
		if st, ok := v.State.(string); ok {
			return st
		}
	}
	return ""
}

// getModeState returns the string value of state of the object named name
// after asserting.
func (m *Viki) getModeState(name string) string {
	if v, ok := m.Objects[name]; ok {
		if st, ok := v.State.(string); ok {
			return st
		}
	}
	return ""
}
