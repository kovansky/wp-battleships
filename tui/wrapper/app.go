package wrapper

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/routines"
	"github.com/kovansky/wp-battleships/ships"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/kovansky/wp-battleships/tui/board"
	"github.com/kovansky/wp-battleships/tui/lobby"
	"github.com/kovansky/wp-battleships/tui/login"
	"github.com/kovansky/wp-battleships/tui/ranking"
	"github.com/kovansky/wp-battleships/tui/setup"
	"github.com/kovansky/wp-battleships/tui/wait"
	"github.com/mbndr/figlet4go"
	"github.com/rs/zerolog"
	"time"
)

type Application struct {
	ctx context.Context
	log zerolog.Logger

	stage tui.Stage
	theme battleships.Theme

	login   login.Login
	lobby   lobby.Lobby
	setup   setup.Setup
	wait    wait.Wait
	game    board.Full
	ranking ranking.Ranking

	width, height int

	asciiRender *figlet4go.AsciiRender
}

func Create(ctx context.Context, theme battleships.Theme) Application {
	asciiRender := figlet4go.NewAsciiRender()

	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	loginApp := login.Create(ctx, theme)

	return Application{
		ctx:         ctx,
		log:         log,
		theme:       theme,
		stage:       tui.StageLogin,
		login:       loginApp,
		asciiRender: asciiRender,
	}
}

func (c Application) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, c.login.Init())

	return tea.Batch(cmds...)
}

func (c Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tmp  tea.Model
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if battleships.Routines.Game != nil {
				battleships.Routines.Game.Quit()
			}
			if battleships.Routines.Lobby != nil {
				battleships.Routines.Lobby.Quit()
			}

			return c, tea.Quit
		}
		break
	case tui.ApplicationStageChangeMsg:
		switch msg.Stage {
		case tui.StageLogin:
			c.stage = msg.Stage

			if msg.From == tui.StageGame {
				battleships.Routines.Game.Quit()
			}
		case tui.StageRanking:
			c.stage = msg.Stage

			players, _ := battleships.ServerClient.Stats()

			if len(battleships.PlayerData.Nick) > 0 {
				player, err := battleships.ServerClient.PlayerStats(battleships.PlayerData.Nick)
				if err == nil {
					players = append(players, ships.NewPlayerFromStats(player))
				}
			}

			c.ranking = ranking.Create(c.ctx, c.theme, players)
			cmds = append(cmds, c.ranking.Init())

			tmp, cmd = c.ranking.Update(tea.WindowSizeMsg{
				Width:  c.width,
				Height: c.height,
			})
			c.ranking = tmp.(ranking.Ranking)
			cmds = append(cmds, cmd)
			break
		case tui.StageLobby:
			c.stage = msg.Stage
			if msg.Model != nil {
				c.lobby = msg.Model.(lobby.Lobby)
			}

			cmds = append(cmds, c.lobby.Init())

			tmp, cmd = c.lobby.Update(tea.WindowSizeMsg{
				Width:  c.width,
				Height: c.height,
			})
			c.lobby = tmp.(lobby.Lobby)
			cmds = append(cmds, cmd)

			battleships.Routines.Lobby = routines.CreateLobby(c.ctx, 3*time.Second, make(chan struct{}))
			go battleships.Routines.Lobby.Run()

			if msg.From == tui.StageGame {
				battleships.Routines.Game.Quit()
			}
			if msg.From == tui.StageWait {
				battleships.Routines.Wait.Quit()
			}
			break
		case tui.StageSetup:
			c.setup = msg.Model.(setup.Setup)
			c.stage = msg.Stage

			cmds = append(cmds, c.setup.Init())

			tmp, cmd = c.setup.Update(tea.WindowSizeMsg{
				Width:  c.width,
				Height: c.height,
			})
			c.setup = tmp.(setup.Setup)
			cmds = append(cmds, cmd)
			break
		case tui.StageWait:
			c.wait = msg.Model.(wait.Wait)
			c.stage = msg.Stage

			cmds = append(cmds, c.wait.Init())

			tmp, cmd = c.wait.Update(tea.WindowSizeMsg{
				Width:  c.width,
				Height: c.height,
			})
			c.wait = tmp.(wait.Wait)
			cmds = append(cmds, cmd)

			battleships.Routines.Wait = routines.CreateWait(c.ctx, 1*time.Second, 7*time.Second, make(chan struct{}))
			go battleships.Routines.Wait.Run()

			break
		case tui.StageGame:
			c.game = msg.Model.(board.Full)
			c.stage = msg.Stage

			cmds = append(cmds, c.game.Init())

			tmp, cmd = c.game.Update(tea.WindowSizeMsg{
				Width:  c.width,
				Height: c.height,
			})
			c.game = tmp.(board.Full)
			cmds = append(cmds, cmd)

			battleships.Routines.Game = routines.CreateGame(c.ctx, 1*time.Second, c.theme, make(chan struct{}))
			go battleships.Routines.Game.Run()

			switch msg.From {
			case tui.StageLobby:
				battleships.Routines.Lobby.Quit()
			case tui.StageWait:
				battleships.Routines.Wait.Quit()
			}
			break
		}
		break
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		break
	}

	switch c.stage {
	case tui.StageLogin:
		tmp, cmd = c.login.Update(msg)
		c.login = tmp.(login.Login)
		cmds = append(cmds, cmd)
		break
	case tui.StageSetup:
		tmp, cmd = c.setup.Update(msg)
		c.setup = tmp.(setup.Setup)
		cmds = append(cmds, cmd)
		break
	case tui.StageWait:
		tmp, cmd = c.wait.Update(msg)
		c.wait = tmp.(wait.Wait)
		cmds = append(cmds, cmd)
		break
	case tui.StageLobby:
		tmp, cmd = c.lobby.Update(msg)
		c.lobby = tmp.(lobby.Lobby)
		cmds = append(cmds, cmd)
		break
	case tui.StageGame:
		tmp, cmd = c.game.Update(msg)
		c.game = tmp.(board.Full)
		cmds = append(cmds, cmd)
		break
	case tui.StageRanking:
		tmp, cmd = c.ranking.Update(msg)
		c.ranking = tmp.(ranking.Ranking)
		cmds = append(cmds, cmd)
		break
	}

	return c, tea.Batch(cmds...)
}

func (c Application) View() string {
	switch c.stage {
	case tui.StageLogin:
		return c.login.View()
	case tui.StageSetup:
		return c.setup.View()
	case tui.StageWait:
		return c.wait.View()
	case tui.StageLobby:
		return c.lobby.View()
	case tui.StageGame:
		return c.game.View()
	case tui.StageRanking:
		return c.ranking.View()
	}

	return ""
}
