package sdk

// NullLogger disregards the log output - and is intended to be used
// when the contents of the debug logger aren't interesting
// to reduce console output
type NullLogger struct{}

// Info prints out a message prefixed with `[INFO]` verbatim
func (NullLogger) Info(_ string) {
}

// Infof prints out a message prefixed with `[INFO]` formatted
// with the specified arguments
func (NullLogger) Infof(_ string, _ ...interface{}) {
}

// Warn prints out a message prefixed with `[WARN]` formatted verbatim
func (NullLogger) Warn(_ string) {
}

// Warnf prints out a message prefixed with `[WARN]` formatted
// with the specified arguments
func (NullLogger) Warnf(_ string, _ ...interface{}) {
}
