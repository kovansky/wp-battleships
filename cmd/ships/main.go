package main

import (
	"context"
	"github.com/kovansky/wp-battleships/ships"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
)

var (
	Version = "v0.0.1"
	client  *ships.Client
	log     zerolog.Logger
)

func main() {
	// Setup signal handlers
	_, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Create log
	log = zerolog.New(os.Stdout).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Create client
	client = ships.NewClient("https://go-pjatk-server.fly.dev/api", &log)

	// Initialize ships
	game, err := client.InitGame(ships.GamePost{
		Desc:  "It's a mee, Mario",
		Nick:  "Mario_the_Plumber",
		Wpbot: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't init the game")
	}

	log.Info().Str("Api-Key", game.Key).Msg("Game started")
}
