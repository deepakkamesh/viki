package device

// Interface implemented by all device types.
type Device interface {
	Execute(interface{}, string)
	On()
	Off()
	Start() error
	Shutdown()
}
