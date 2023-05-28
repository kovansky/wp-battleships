package wrapper

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/routines"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/kovansky/wp-battleships/tui/board"
	"github.com/kovansky/wp-battleships/tui/lobby"
	"github.com/kovansky/wp-battleships/tui/login"
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

	login login.Login
	lobby lobby.Lobby
	wait  wait.Wait
	game  board.Full

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
		case "ctrl+c", "q":
			if battleships.Routines.Game != nil {
				battleships.Routines.Game.Quit()
			}
			if battleships.Routines.Lobby != nil {
				battleships.Routines.Lobby.Quit()
			}

			return c, tea.Quit
		}
	case tui.ApplicationStageChangeMsg:
		switch msg.Stage {
		case tui.StageLobby:
			c.lobby = msg.Model.(lobby.Lobby)
			c.stage = msg.Stage

			cmds = append(cmds, c.lobby.Init())

			tmp, cmd = c.lobby.Update(tea.WindowSizeMsg{
				Width:  c.width,
				Height: c.height,
			})
			c.lobby = tmp.(lobby.Lobby)
			cmds = append(cmds, cmd)

			battleships.Routines.Lobby = routines.CreateLobby(c.ctx, 5*time.Second, make(chan struct{}))
			go battleships.Routines.Lobby.Run()

			if msg.From == tui.StageGame {
				battleships.Routines.Game.Quit()
			}
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
			battleships.Routines.Lobby.Quit()
		}
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
	}

	if c.stage == tui.StageLogin {
		tmp, cmd = c.login.Update(msg)
		c.login = tmp.(login.Login)
		cmds = append(cmds, cmd)
	}

	if c.stage == tui.StageWait {
		tmp, cmd = c.wait.Update(msg)
		c.wait = tmp.(wait.Wait)
		cmds = append(cmds, cmd)
	}

	if c.stage == tui.StageLobby {
		tmp, cmd = c.lobby.Update(msg)
		c.lobby = tmp.(lobby.Lobby)
		cmds = append(cmds, cmd)
	}

	if c.stage == tui.StageGame {
		tmp, cmd = c.game.Update(msg)
		c.game = tmp.(board.Full)
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

func (c Application) View() string {
	switch c.stage {
	case tui.StageLogin:
		return c.login.View()
	case tui.StageWait:
		return c.wait.View()
	case tui.StageLobby:
		return c.lobby.View()
	case tui.StageGame:
		return c.game.View()
	}

	return ""
}
