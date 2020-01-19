package cookiemonster

// Logger represents the log output of the worker pool
type Logger interface {
	Printf(format string, args ...interface{})
}
