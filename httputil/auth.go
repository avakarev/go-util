package httputil

import (
	"strings"
)

// RequestContext represents http request context
type RequestContext interface {
	GetHeader(s string) string
}

// AuthBearer extracts bearer token value from Authorization header
func AuthBearer(ctx RequestContext) string {
	t := ctx.GetHeader("Authorization")
	if len(t) > 6 && strings.ToLower(t[0:6]) == "bearer" {
		return strings.TrimSpace(t[7:])
	}
	return ""
}
