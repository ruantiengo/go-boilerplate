package logger

import (
	"os"

	"github.com/fatih/color"
)

type LogType int

const (
	Info LogType = iota
	Warning
	Error
	Debug
	Default
)

func Message(logType LogType, message string, args ...interface{}) {
	switch logType {
	case Info:
		color.Green(message, args...)
	case Warning:
		color.Yellow(message, args...)
	case Error:
		color.Red(message, args...)
	case Debug:
		debug := os.Getenv("DEBUG_MODE") == "true"
		if debug {
			color.Blue(message, args...)
		}
	default:
		color.White(message, args...)
	}
}
