// Package procutil provides util helpers to work with processes and signals
package procutil

import (
	"os"
	"os/signal"
)

// WaitSignal waits for given signals
// Defautls to `os.Interrupt` if given no signals
func WaitSignal(sigs ...os.Signal) os.Signal {
	if len(sigs) == 0 { // default signal
		sigs = append(sigs, os.Interrupt)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)
	return <-ch
}
