package cookiemonster

type Logger interface {
	Printf(format string, args ...interface{})
}
