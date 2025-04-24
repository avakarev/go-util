// Package wsutil implements websocket communication
package wsutil

// Event defines generic ws event
type Event struct {
	Topic string
	Data  []byte
}
