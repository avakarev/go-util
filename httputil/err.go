package httputil

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avakarev/go-util/strutil"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

// ValidationErr represents validation error
type ValidationErr struct {
	Subject string `json:"subject"`
	Msg     string `json:"msg"`
}

// Err represents generic api error
type Err struct {
	Code  int             `json:"code"`
	Msg   string          `json:"msg"`
	Items []ValidationErr `json:"items,omitempty"`
}

// ErrResponse represents json container for error object
type ErrResponse struct {
	Error Err `json:"error"`
}

// WriteJSON responds json error with http.ResponseWriter
func (e *ErrResponse) WriteJSON(w http.ResponseWriter) {
	w.WriteHeader(e.Error.Code)
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Send()
	}
	if _, err := w.Write([]byte(bytes)); err != nil {
		log.Error().Err(err).Send()
	}
}

// StdErrMsg returns standard error message by code
func StdErrMsg(code int) string {
	switch code {
	case http.StatusBadRequest:
		return "bad request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not found"
	case http.StatusInternalServerError:
		return "internal server error"
	default:
		return "unknown error"
	}
}

// ValidationErrMsg returns validation error message by validate tag
func ValidationErrMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "required"
	case "email":
		return "invalid email format"
	}
	return "invalid"
}

// NewErr returns new standard api error value
func NewErr(code int, msg string) *ErrResponse {
	if msg == "" {
		msg = StdErrMsg(code)
	}
	return &ErrResponse{
		Error: Err{
			Code: code,
			Msg:  msg,
		},
	}
}

// NewValidationErr returns new validation error value
func NewValidationErr(errors validator.ValidationErrors) *ErrResponse {
	err := NewErr(http.StatusBadRequest, "validation error")
	err.Error.Items = make([]ValidationErr, len(errors))
	for i, e := range errors {
		err.Error.Items[i] = ValidationErr{
			Subject: strutil.Decapitalize(e.Field()),
			Msg:     ValidationErrMsg(e),
		}
	}
	return err
}

// NewErrFrom returns new error value from given error
func NewErrFrom(err error) *ErrResponse {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return NewValidationErr(ve)
	}
	return NewErr(http.StatusInternalServerError, err.Error())
}
