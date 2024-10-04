package httputil

import (
	"strings"
)

// AuthBearer extracts bearer token value from Authorization header
func AuthBearer(auth string) string {
	if len(auth) > 6 && strings.ToLower(auth[0:6]) == "bearer" {
		return strings.TrimSpace(auth[7:])
	}
	return ""
}
