package timeutil

// TickFn defines a function that timer calls
type TickFn func()

// Timer defines generic timer interface
type Timer interface {
	Start()
	Stop()
}
