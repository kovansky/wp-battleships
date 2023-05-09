package main

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/ships"
	"github.com/kovansky/wp-battleships/tui"
	board2 "github.com/kovansky/wp-battleships/tui/board"
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
		Wpbot: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't init the game")
	}

	log.Info().Str("Api-Key", game.Key()).Msg("Game started")

	err = battleships.ServerClient.UpdateBoard(game)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't update the game board")
	}
	err = battleships.ServerClient.GameStatus(game)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't update the game status")
	}

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
				Foreground(lipgloss.Color("#ff0000"))))
	globalTheme := tui.NewTheme().
		SetTextPrimary(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd700"))).
		SetTextSecondary(lipgloss.NewStyle().Foreground(lipgloss.Color("#1e90ff")))

	playersInfo := fmt.Sprintf(lipgloss.NewStyle().Italic(true).Render("Waiting for game..."))

	boardComponent := board2.InitFull(game, theme, theme, globalTheme, playersInfo)

	program := tea.NewProgram(boardComponent, tea.WithAltScreen())

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err = battleships.ServerClient.GameStatus(game)
				if err != nil {
					log.Fatal().Err(err).Msg("Couldn't update the game status")
				}

				if game.GameStatus().Status == battleships.StatusGameInProgress && game.GameStatus().ShouldFire && (game.Opponent() == nil || game.Opponent().Name() == "") {
					err = battleships.ServerClient.GameDesc(game)
					if err != nil {
						log.Fatal().Err(err).Msg("Couldn't update the game description")
					}

					playersInfo = fmt.Sprintf("%s %s %s\n"+
						"%s %s %s",
						globalTheme.TextPrimary.Copy().Bold(true).Render("YOU"),
						game.Player().Name(),
						lipgloss.NewStyle().Italic(true).Render("("+game.Player().Description()+")"),
						globalTheme.TextPrimary.Copy().Bold(true).Render("ENEMY"),
						game.Opponent().Name(),
						lipgloss.NewStyle().Italic(true).Render("("+game.Opponent().Description()+")"),
					)

					program.Send(battleships.PlayersUpdateMsg{PlayersInfo: playersInfo})
				}

				program.Send(battleships.GameUpdateMsg{})
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	if _, err := program.Run(); err != nil {
		close(quit)
		log.Error().Err(err).Msg("Could not draw board")
	}
}
