// Package natsutil implements nats helpers
package natsutil

import (
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/avakarev/go-util/httputil"
)

// Respond responds given bytes
func Respond(msg *nats.Msg, bytes []byte) error {
	return msg.Respond(bytes)
}

// RespondJSON responds given value as marshalled bytes
func RespondJSON(msg *nats.Msg, v any) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return Respond(msg, bytes)
}

// RespondJSONErr responds given error value as marshalled bytes
func RespondJSONErr(msg *nats.Msg, err error) error {
	return RespondJSON(msg, httputil.NewErrFrom(err))
}
