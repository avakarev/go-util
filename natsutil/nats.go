package natsutil

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/avakarev/go-util/httputil"
)

// RespondJSON responds given value as marshalled bytes
func RespondJSON(msg *nats.Msg, v any) {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	if err := msg.Respond(bytes); err != nil {
		log.Error().Err(err).Send()
	}
}

// RespondJSONErr responds given error value as marshalled bytes
func RespondJSONErr(msg *nats.Msg, err error) {
	RespondJSON(msg, httputil.NewErrFrom(err))
}
