package sdk

// Logger is an interface for switching out the Logger implementation
type Logger interface {
	// Info prints out a message prefixed with `[INFO]` verbatim
	Info(message string)

	// Infof prints out a message prefixed with `[INFO]` formatted
	// with the specified arguments
	Infof(format string, args ...interface{})

	// Warn prints out a message prefixed with `[WARN]` formatted verbatim
	Warn(message string)

	// Warnf prints out a message prefixed with `[WARN]` formatted
	// with the specified arguments
	Warnf(format string, args ...interface{})
}
