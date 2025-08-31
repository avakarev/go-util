package natsutil

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

// MsgHandlerFunc defines callback function that processes messages delivered to asynchronous subscribers
type MsgHandlerFunc func(msg *nats.Msg) error

// ErrHandlerFunc defines error handlers invoked in case when MsgHandlerFunc returns error
type ErrHandlerFunc func(msg *nats.Msg, err error)

// DefaultErrHandler implements default error handler
// It just logs error altogether with message subject
func DefaultErrHandler(msg *nats.Msg, err error) {
	log.Error().Err(err).Str("subject", msg.Subject).Msg("nats: error handler")
	// if there is no reply field, this is simple published message, nothing to do
	if msg.Reply == "" {
		return
	}
	// otherwise this is request reply, try to reply with json-serialized error
	if err := RespondJSONErr(msg, err); err != nil {
		log.Error().Err(err).Str("subject", msg.Subject).Msg("nats: error handler failed")
	}
}
