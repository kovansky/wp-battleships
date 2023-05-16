package wrapper

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/routines"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/kovansky/wp-battleships/tui/board"
	"github.com/kovansky/wp-battleships/tui/lobby"
	"github.com/mbndr/figlet4go"
	"github.com/rs/zerolog"
	"time"
)

type Application struct {
	ctx context.Context
	log zerolog.Logger

	stage tui.Stage
	theme battleships.Theme

	lobby lobby.Lobby
	game  board.Full

	width, height int

	asciiRender *figlet4go.AsciiRender

	quit chan struct{}
}

func Create(ctx context.Context, theme battleships.Theme, lobbyApp lobby.Lobby) Application {
	asciiRender := figlet4go.NewAsciiRender()

	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	return Application{
		ctx:         ctx,
		log:         log,
		theme:       theme,
		stage:       tui.StageLobby,
		lobby:       lobbyApp,
		asciiRender: asciiRender,
	}
}

func (c Application) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, c.lobby.Init())

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
			close(c.quit)
			battleships.Routines.Lobby.Quit()

			return c, tea.Quit
		}
	case tui.ApplicationStageChangeMsg:
		switch msg.Stage {
		case tui.StageGame:
			c.game = msg.Board.(board.Full)
			c.stage = msg.Stage

			cmds = append(cmds, c.game.Init())

			tmp, cmd = c.game.Update(tea.WindowSizeMsg{
				Width:  c.width,
				Height: c.height,
			})
			c.game = tmp.(board.Full)
			cmds = append(cmds, cmd)

			c.quit = make(chan struct{})
			gameRoutine := routines.CreateGame(c.ctx, 1*time.Second, c.theme, c.quit)
			go gameRoutine.Run()
			battleships.Routines.Lobby.Quit()
		}
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
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
	case tui.StageLobby:
		return c.lobby.View()
	case tui.StageGame:
		return c.game.View()
	}

	return ""
}
