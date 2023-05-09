package main

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/ships"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/kovansky/wp-battleships/tui/lobby"
	"github.com/rs/zerolog"
	"os"
)

var (
	Version = "v0.0.1"
	log     zerolog.Logger
)

func main() {
	// Propagate build info
	battleships.Version = Version

	// Setup signal handlers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create logger
	log = zerolog.
		New(os.Stdout).
		With().Timestamp().
		Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Create client
	battleships.ServerClient = ships.NewClient(ctx, "https://go-pjatk-server.fly.dev/api", &log)

	// Initialize ships
	game, err := battleships.ServerClient.InitGame(battleships.GamePost{
		Desc:  "It's a mee, Mario",
		Nick:  "Mario_the_Plumber",
		Wpbot: false,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't init the game")
	}

	log.Info().Str("Api-Key", game.Key()).Msg("Game started")

	err = battleships.ServerClient.GameStatus(game)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't update the game status")
	}

	globalTheme := tui.NewTheme().
		SetTextPrimary(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd700"))).
		SetTextSecondary(lipgloss.NewStyle().Foreground(lipgloss.Color("#1e90ff")))

	lobbyComponent := lobby.Create(globalTheme)

	program := tea.NewProgram(lobbyComponent, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Error().Err(err).Msg("Could not draw board")
	}
}
