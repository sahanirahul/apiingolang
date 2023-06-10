package logging

type Level int8

// this is a copy of zap level, can use zap.Level instead
const (
	DebugLevel Level = iota - 1

	InfoLevel

	WarnLevel

	ErrorLevel

	PanicLevel
)
