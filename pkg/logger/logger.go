package logger

import (
	"io"
	"log"
	"os"
)

var (
	// Info logger
	Info *log.Logger
	// Debug loggger
	Debug *log.Logger
)

func init() {
	if Info == nil {
		Info = log.New(io.Discard, "INFO: XDPFail2Ban: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	if Debug == nil {
		Debug = log.New(io.Discard, "DEBUG: XDPFail2Ban: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

// SetLevel of the log. Options: INFO|DEBUG
func SetLevel(logLevel string) {
	switch logLevel {
	case "INFO":
		Info.SetOutput(os.Stdout)
	case "DEBUG":
		Info.SetOutput(os.Stdout)
		Debug.SetOutput(os.Stdout)
	}
}
