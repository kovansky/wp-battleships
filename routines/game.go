package routines

import (
	"context"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/rs/zerolog"
	"time"
)

type Game struct {
	log zerolog.Logger

	theme battleships.Theme

	duration time.Duration

	quit chan struct{}
}

func CreateGame(ctx context.Context, duration time.Duration, theme battleships.Theme, quit chan struct{}) Game {
	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	return Game{
		log:      log,
		duration: duration,
		quit:     quit,
		theme:    theme,
	}
}

func (g Game) Run() {
	ticker := time.NewTicker(g.duration)
	for {
		select {
		case <-ticker.C:
			err := battleships.ServerClient.GameStatus(battleships.GameInstance)
			if err != nil {
				g.log.Fatal().Err(err).Msg("Couldn't update the game status")
			}

			if battleships.GameInstance.GameStatus().Status == battleships.StatusGameInProgress && battleships.GameInstance.GameStatus().ShouldFire && (battleships.GameInstance.Opponent() == nil || battleships.GameInstance.Opponent().Name() == "") {
				err = battleships.ServerClient.GameDesc(battleships.GameInstance)
				if err != nil {
					g.log.Fatal().Err(err).Msg("Couldn't update the game description")
				}

				playersInfo := fmt.Sprintf("%s %s %s\n"+
					"%s %s %s",
					g.theme.TextPrimary().Copy().Bold(true).Render("YOU"),
					battleships.GameInstance.Player().Name(),
					lipgloss.NewStyle().Italic(true).Render("("+battleships.GameInstance.Player().Description()+")"),
					g.theme.TextPrimary().Copy().Bold(true).Render("ENEMY"),
					battleships.GameInstance.Opponent().Name(),
					lipgloss.NewStyle().Italic(true).Render("("+battleships.GameInstance.Opponent().Description()+")"),
				)

				battleships.ProgramMessage(battleships.PlayersUpdateMsg{PlayersInfo: playersInfo})
			}

			battleships.ProgramMessage(battleships.GameUpdateMsg{})
		case <-g.quit:
			ticker.Stop()
			return
		}
	}
}

func (g Game) Quit() {
	select {
	case <-g.quit:
	default:
		if battleships.GameInstance != nil && battleships.GameInstance.Key() != "" {
			_ = battleships.ServerClient.Abandon(battleships.GameInstance)
		}

		close(g.quit)
	}
}
