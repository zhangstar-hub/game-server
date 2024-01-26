package throttle

type RequestTrottle interface {
	Close()
	CanRequest() bool
}
