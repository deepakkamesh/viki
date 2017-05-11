package viki

// getMochadState returns the string value of state of the object named name
// after asserting.
func (m *Viki) getMochadState(name string) string {
	_, o := m.ObjectManager.GetObjectByName(name)
	if o == nil {
		return ""
	}
	st, ok := o.State.(string)
	if !ok {
		return ""
	}
	return st
}

// getModeState returns the string value of state of the object named name
// after asserting.
func (m *Viki) getModeState(name string) string {
	_, o := m.ObjectManager.GetObjectByName(name)
	if o == nil {
		return ""
	}
	st, ok := o.State.(string)
	if !ok {
		return ""
	}
	return st

}
