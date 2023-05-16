package routines

import (
	"context"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/rs/zerolog"
	"time"
)

type Lobby struct {
	log zerolog.Logger

	duration time.Duration

	quit chan struct{}
}

func CreateLobby(ctx context.Context, duration time.Duration, quit chan struct{}) Lobby {
	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	return Lobby{
		log:      log,
		duration: duration,
		quit:     quit,
	}
}

func (l Lobby) Run() {
	ticker := time.NewTicker(l.duration)
	for {
		select {
		case <-ticker.C:
			players, err := battleships.ServerClient.ListPlayers()
			if err != nil {
				l.log.Fatal().Err(err).Msg("Couldn't list players")
			}

			battleships.ProgramMessage(battleships.PlayersListMsg{Players: players})
		case <-l.quit:
			ticker.Stop()
			return
		}
	}
}

func (l Lobby) Quit() {
	select {
	case <-l.quit:
	default:
		close(l.quit)
	}
}
