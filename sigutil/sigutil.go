// Package sigutil provides signal helpers
package sigutil

import (
	"os"
	"os/signal"
	"syscall"
)

// Wait waits for first of given signals
// Defaults to {syscall.SIGINT, syscall.SIGTERM} if no signals given
func Wait(sigs ...os.Signal) os.Signal {
	if len(sigs) == 0 { // default
		sigs = append(sigs, syscall.SIGINT, syscall.SIGTERM)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)
	return <-ch
}
