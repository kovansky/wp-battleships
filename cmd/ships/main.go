package main

import (
	tea "github.com/charmbracelet/bubbletea"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/board"
	"github.com/kovansky/wp-battleships/ships"
	"github.com/rs/zerolog"
	"os"
)

var (
	Version = "v0.0.1"
	client  *ships.Client
	log     zerolog.Logger
)

func main() {
	// Create log
	log = zerolog.New(os.Stdout).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Create client
	client = ships.NewClient("https://go-pjatk-server.fly.dev/api", &log)

	// Initialize ships
	game, err := client.InitGame(battleships.GamePost{
		Desc:  "It's a mee, Mario",
		Nick:  "Mario_the_Plumber",
		Wpbot: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't init the game")
	}

	log.Info().Str("Api-Key", game.Key()).Msg("Game started")

	p := tea.NewProgram(board.InitialBoard([]string{"A1", "A10", "B10"}))
	if _, err = p.Run(); err != nil {
		log.Fatal().Err(err).Msg("Could not run the GUI")
	}
}
