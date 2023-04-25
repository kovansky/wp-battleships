package main

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/board"
	"github.com/kovansky/wp-battleships/ships"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"time"
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
	err = client.GameDesc(game)
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
	enemyTheme := board.NewTheme().
		SetRows(colRowStyle).
		SetCols(colRowStyle).
		SetShip(board.NewBrush().
			SetChar('X').
			SetStyle(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff4500"))))
	globalTheme := board.NewTheme().
		SetTextPrimary(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd700"))).
		SetTextSecondary(lipgloss.NewStyle().Foreground(lipgloss.Color("#1e90ff")))

	playersInfo := fmt.Sprintf("%s %s %s\n"+
		"%s %s %s",
		globalTheme.TextPrimary.Copy().Bold(true).Render("YOU"),
		game.Player().Name(),
		lipgloss.NewStyle().Italic(true).Render("("+game.Player().Description()+")"),
		globalTheme.TextPrimary.Copy().Bold(true).Render("ENEMY"),
		game.Opponent().Name(),
		lipgloss.NewStyle().Italic(true).Render("("+game.Opponent().Description()+")"),
	)

	boardComponent := board.InitFullComponent(theme, enemyTheme, globalTheme, playersInfo, game.Board(), []string{})

	program := tea.NewProgram(boardComponent, tea.WithAltScreen())

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err = client.GameStatus(game)
				if err != nil {
					log.Fatal().Err(err).Msg("Couldn't update the game status")
				}
				//log.Debug().Interface("gameStatus", game.GameStatus()).Msg("game status updated")
				if game.GameStatus().Status == battleships.StatusGameInProgress && game.Opponent().Name() == "" {
					err = client.GameDesc(game)
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

				program.Send(battleships.GameUpdateMsg{GameStatus: game.GameStatus()})
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
