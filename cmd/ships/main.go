package main

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/board"
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
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Create log
	log = zerolog.New(os.Stdout).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Create client
	client = ships.NewClient(ctx, "https://go-pjatk-server.fly.dev/api", &log)

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

	err = client.UpdateBoard(game)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't update the game board")
	}
	err = client.GameStatus(game)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't update the game status")
	}

	colRowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff7f")).
		Bold(true)
	theme := board.NewTheme().
		SetRows(colRowStyle).
		SetCols(colRowStyle).
		SetShip(board.NewBrush().
			SetChar('X').
			SetStyle(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff4500"))))

	program := tea.NewProgram(board.InitComponent(theme, game.Board()...))
	if _, err := program.Run(); err != nil {
		log.Error().Err(err).Msg("Could not draw board")
	}

	//log.Info().
	//	Str("Name", game.Player().Name()).
	//	Str("Description", game.Player().Description()).
	//	Msg("You")
	//log.Info().
	//	Str("Name", game.Opponent().Name()).
	//	Str("Description", game.Opponent().Description()).
	//	Msg("Opponent")
}
