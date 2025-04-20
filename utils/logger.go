package utils

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func NewLogger() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
}
