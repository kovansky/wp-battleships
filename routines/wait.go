package routines

import (
	"context"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/kovansky/wp-battleships/tui/board"
	"github.com/rs/zerolog"
	"time"
)

type Wait struct {
	log zerolog.Logger

	statusDuration  time.Duration
	refreshDuration time.Duration
	theme           battleships.Theme

	quit chan struct{}
}

func CreateWait(ctx context.Context, statusDuration time.Duration, refreshDuration time.Duration, quit chan struct{}) Wait {
	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	return Wait{
		log:             log,
		statusDuration:  statusDuration,
		refreshDuration: refreshDuration,
		theme:           battleships.Themes.Global,
		quit:            quit,
	}
}

func (w Wait) Run() {
	statusTicker := time.NewTicker(w.statusDuration)
	refreshTicker := time.NewTicker(w.refreshDuration)
	for {
		select {
		case <-statusTicker.C:
			err := battleships.ServerClient.GameStatus(battleships.GameInstance)
			if err != nil {
				w.log.Fatal().Err(err).Msg("Could not update game status")
			}

			if battleships.GameInstance.GameStatus().Status == battleships.StatusGameInProgress {
				err = battleships.ServerClient.UpdateBoard(battleships.GameInstance)
				if err != nil {
					w.log.Fatal().Err(err).Msg("Couldn't update the game board")
				}
				err = battleships.ServerClient.GameDesc(battleships.GameInstance)
				if err != nil {
					w.log.Fatal().Err(err).Msg("Couldn't update the game status")
				}

				playersInfo := fmt.Sprintf("%s %s %s\n"+
					"%s %s %s",
					w.theme.TextPrimary().Copy().Bold(true).Render("YOU"),
					battleships.GameInstance.Player().Name(),
					lipgloss.NewStyle().Italic(true).Render("("+battleships.GameInstance.Player().Description()+")"),
					w.theme.TextPrimary().Copy().Bold(true).Render("ENEMY"),
					battleships.GameInstance.Opponent().Name(),
					lipgloss.NewStyle().Italic(true).Render("("+battleships.GameInstance.Opponent().Description()+")"),
				)

				gameBoard := board.InitFull(battleships.GameInstance, battleships.Themes.Player, battleships.Themes.Enemy, battleships.Themes.Global, playersInfo)

				battleships.ProgramMessage(tui.ApplicationStageChangeMsg{
					From:  tui.StageWait,
					Stage: tui.StageGame,
					Model: gameBoard,
				})
			}
		case <-refreshTicker.C:
			err := battleships.ServerClient.Refresh(battleships.GameInstance)
			if err != nil {
				statusErr := battleships.ServerClient.GameStatus(battleships.GameInstance)
				if statusErr != nil {
					w.log.Fatal().Err(statusErr).Msg("Could not update game status")
				}

				if battleships.GameInstance.GameStatus().Status != battleships.StatusGameInProgress {
					w.log.Fatal().Err(err).Msg("Couldn't refresh game")
				}
			}
		case <-w.quit:
			statusTicker.Stop()
			refreshTicker.Stop()
			return
		}
	}
}

func (w Wait) Quit() {
	select {
	case <-w.quit:
	default:
		close(w.quit)
	}
}
