package wsutil

// ErrorHandler defines a function that will process all errors
type ErrorHandler = func(error)

// DefaultErrorHandler that process return errors from handlers
func DefaultErrorHandler(_ error) {}

// Config is a struct holding the hub settings
type Config struct {
	// ErrorHandler is executed on error
	ErrorHandler ErrorHandler
}
