package cookiejar

type Logger interface {
	Printf(format string, args ...interface{})
}
