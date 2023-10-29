package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

type logger struct {
	log zerolog.Logger
}

func New(logFile string) (Logger, error) {
	file, err := os.Create(logFile)
	if err != nil {
		return nil, err
	}

	zerolog.CallerSkipFrameCount = 4

	consoleWriter := zerolog.ConsoleWriter{
		Out:     file,
		NoColor: true,
	}

	zeroLogger := zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()

	return logger{
		log: zeroLogger,
	}, nil
}

func (l logger) Info(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.log.Info().Msg(msg)
		return
	}

	l.log.Info().Msg(fmt.Sprintf(msg, args...))
}

func (l logger) Debug(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.log.Debug().Msg(msg)
		return
	}

	l.log.Debug().Msg(fmt.Sprintf(msg, args...))
}

func (l logger) Error(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.log.Error().Msg(msg)
		return
	}

	l.log.Error().Msg(fmt.Sprintf(msg, args...))
}

func (l logger) Fatal(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.log.Fatal().Msg(msg)
		return
	}

	l.log.Fatal().Msg(fmt.Sprintf(msg, args...))
}
