package appcore

type ILogger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Fatal(message string)
	Flush() error
}
