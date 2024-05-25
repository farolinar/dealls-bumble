package main

import (
	"os"
	"time"

	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/credentials"
	// "github.com/aws/aws-sdk-go/aws/session"
	"github.com/farolinar/dealls-bumble/cmd/app"
	"github.com/farolinar/dealls-bumble/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	envConfig, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error loading env config")
	}

	log.Logger = log.With().Caller().Logger()
	zerolog.TimeFieldFormat = time.RFC3339

	if envConfig.App.LogPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, FormatTimestamp: func(i interface{}) string { return time.Now().Format(time.RFC3339) }})
	}

	switch envConfig.App.LogLevel {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	app.Serve()
}
