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
	"time"
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
	game, err := battleships.ServerClient.InitGame(battleships.GamePost{
		Wpbot: false,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't init the game")
	}

	log.Info().Str("Api-Key", game.Key()).Msg("Game started")

	err = battleships.ServerClient.GameDesc(game)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't update the game status")
	}
	players, err := battleships.ServerClient.ListPlayers()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't list players")
	}
	for i, player := range players {
		if player.Name() == game.Player().Name() {
			players = removeFromSlice(players, i)
			break
		}
	}

	globalTheme := tui.NewTheme().
		SetTextPrimary(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd700"))).
		SetTextSecondary(lipgloss.NewStyle().Foreground(lipgloss.Color("#1e90ff")))

	lobbyComponent := lobby.Create(ctx, globalTheme, players)

	program := tea.NewProgram(lobbyComponent, tea.WithAltScreen())

	ticker := time.NewTicker(5 * time.Second)
	refreshTicker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				players, err := battleships.ServerClient.ListPlayers()
				if err != nil {
					log.Fatal().Err(err).Msg("Couldn't list players")
				}

				program.Send(battleships.PlayersListMsg{Players: players})
			case <-refreshTicker.C:
				err := battleships.ServerClient.Refresh(game)
				if err != nil {
					log.Fatal().Err(err).Msg("Couldn't list players")
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	if _, err := program.Run(); err != nil {
		log.Error().Err(err).Msg("Could not draw board")
	}
}

func removeFromSlice[T any](s []T, i int) []T {
	if i < 0 || i >= len(s) {
		return s
	}

	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
