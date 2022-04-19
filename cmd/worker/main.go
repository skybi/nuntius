package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/skybi/nuntius/internal/client"
	"github.com/skybi/nuntius/internal/config"
	"github.com/skybi/nuntius/internal/metar"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Set up zerolog to use pretty printing
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})
	log.Info().Msg("starting up...")

	// Load the application configuration
	log.Info().Msg("loading configuration...")
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("could not load the configuration")
	}
	if cfg.IsEnvProduction() {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Debug().Str("config", fmt.Sprintf("%+v", cfg)).Msg("")

	// Initialize the API client
	apiClient := client.New(cfg.APIAddress, cfg.APIKey)
	keyInfo, err := apiClient.GetKeyInfo()
	if err != nil {
		log.Fatal().Err(err).Msg("could not retrieve API key information")
	}

	// Check if the API key has unlimited quota and rate limit
	if keyInfo.Quota >= 0 {
		log.Fatal().Int64("quota", keyInfo.Quota).Msg("an API key with unlimited quota is required")
	}
	if keyInfo.RateLimit >= 0 {
		log.Fatal().Int("rate_limit", keyInfo.RateLimit).Msg("an API key with no rate limit is required")
	}

	// Check if the key may feed METARs
	if cfg.FeedMETARs && keyInfo.Capabilities&client.CapabilityFeedMETARs == 0 {
		cfg.FeedMETARs = false
		log.Warn().Msg("METAR feeding disabled due to lack of required key capability")
	}

	// Abort if no feeding is enabled
	if !cfg.FeedMETARs {
		log.Fatal().Msg("aborting due to disabled feeding")
	}

	// Start feeding METARs if necessary
	if cfg.FeedMETARs {
		log.Info().Msg("starting the METAR feeder...")
		feeder := metar.NewFeeder(apiClient, 100, time.Second)
		feeder.Start()
		defer feeder.Stop()
	}

	// TODO: startup logic

	// Wait for the application to be terminated
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
}
