package natsutil

// ConfigureFunc defines configurator func
type ConfigureFunc func(*Conn) error

// WithLatencyLog enables latency logging
func WithLatencyLog(ignoreSubjs ...string) ConfigureFunc {
	return func(conn *Conn) error {
		conn.latencyLog.enabled = true
		for _, subj := range ignoreSubjs {
			conn.latencyLog.ignoreSubjects[conn.enrichSubj(subj)] = true
		}
		return nil
	}
}
