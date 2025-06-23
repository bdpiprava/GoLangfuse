package types

// Level observation log level, supported levels are Debug, Default, Warning and Error
type Level string

const (
	Debug   Level = "DEBUG"   // Debug for logging debug logs
	Default Level = "DEFAULT" // Default for logging default logs
	Warning Level = "WARNING" // Warning for logging warning logs
	Error   Level = "ERROR"   // Error for logging error logs
)
