// Package logging is to define our logger
package logging

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	log Logger
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	output := defaultZerologConsoleWriter()

	if os.Getenv("LOG_MODE") == "debug" {
		log = Logger{
			zerolog.New(output).With().Timestamp().Caller().Logger(),
		}
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		log = Logger{
			zerolog.New(output).With().Timestamp().Logger(),
		}
	}
}

// Logger is the wrapper of zerolog.Logger
type Logger struct {
	zerolog.Logger
}

func (l *Logger) Msg() {

}

func defaultZerologConsoleWriter() zerolog.ConsoleWriter {
	return zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    true,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			level := strings.ToUpper(fmt.Sprintf("%s:", i))
			switch level {
			case "DEBUG:":
				level = Yellow(level)
			case "ERROR:":
				level = Red(level)
			}

			return level
		},
	}
}

// GetLogger to return the global log.
func GetLogger() *Logger {
	return &log
}
