package main

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/ships"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/kovansky/wp-battleships/tui/wrapper"
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

	// Create logger
	log = zerolog.
		New(os.Stdout).
		With().Timestamp().
		Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Setup signal handlers
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, battleships.ContextKeyLog, log)
	defer cancel()

	// Create client
	battleships.ServerClient = ships.NewClient(ctx, "https://go-pjatk-server.fly.dev/api", &log)

	// Initialize ships
	colRowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff7f")).
		Bold(true)
	theme := tui.NewTheme().
		SetRows(colRowStyle).
		SetCols(colRowStyle).
		SetShip(tui.NewBrush().
			SetChar('X').
			SetStyle(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1e90ff")))).
		SetHit(tui.NewBrush().
			SetChar('X').
			SetStyle(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ff7f")))).
		SetSunk(tui.NewBrush().
			SetChar('-').
			SetStyle(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#006332")))).
		SetMiss(tui.NewBrush().
			SetChar('o').
			SetStyle(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff0000")))).
		SetPotential(tui.NewBrush().
			SetChar('o').
			SetStyle(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e06c00"))))
	globalTheme := tui.NewTheme().
		SetTextPrimary(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd700"))).
		SetTextSecondary(lipgloss.NewStyle().Foreground(lipgloss.Color("#1e90ff")))

	battleships.Themes = battleships.GameThemes{
		Player: theme,
		Enemy:  theme,
		Global: globalTheme,
	}

	applicationWrapper := wrapper.Create(ctx, globalTheme)

	program := tea.NewProgram(applicationWrapper, tea.WithAltScreen())

	battleships.ProgramMessage = func(msg tea.Msg) {
		program.Send(msg)
	}

	if _, err := program.Run(); err != nil {
		log.Error().Err(err).Msg("Could not draw board")
	}
}
