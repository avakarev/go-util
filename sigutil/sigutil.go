// Package sigutil provides signal helpers
package sigutil

import (
	"os"
	"os/signal"
)

// Wait waits for first of given signals
// Defautls to `os.Interrupt` if no signals given
func Wait(sigs ...os.Signal) os.Signal {
	if len(sigs) == 0 { // default signal
		sigs = append(sigs, os.Interrupt)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)
	return <-ch
}
