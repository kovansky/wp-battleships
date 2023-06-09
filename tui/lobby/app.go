package lobby

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/kovansky/wp-battleships/tui/common"
	"github.com/mbndr/figlet4go"
	"github.com/rs/zerolog"
)

type Lobby struct {
	log zerolog.Logger

	theme battleships.Theme

	subcomponents  map[string]tea.Model
	initialPlayers []battleships.Player

	asciiRender *figlet4go.AsciiRender
}

func Create(ctx context.Context, theme battleships.Theme, initialPlayers []battleships.Player) Lobby {
	asciiRender := figlet4go.NewAsciiRender()
	header := common.CreateHeader("Battleships", theme, asciiRender)
	table := CreatePlayers(ctx, theme, initialPlayers)

	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	return Lobby{
		log:         log,
		theme:       theme,
		asciiRender: asciiRender,
		subcomponents: map[string]tea.Model{
			"header": header,
			"table":  table,
		},
		initialPlayers: initialPlayers,
	}
}

func (c Lobby) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, cmp := range c.subcomponents {
		cmds = append(cmds, cmp.Init())
	}

	return tea.Batch(cmds...)
}

func (c Lobby) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return c, tea.Quit
		case "R":
			return c, func() tea.Msg {
				return tui.ApplicationStageChangeMsg{
					From:  tui.StageLogin,
					Stage: tui.StageRanking,
				}
			}
		}
	case tea.WindowSizeMsg:
		for name, cmp := range c.subcomponents {
			c.subcomponents[name], cmd = cmp.Update(msg)
			cmds = append(cmds, cmd)
		}

		table := c.subcomponents["table"].(Players)

		table.Width = msg.Width
		table.Height = msg.Height - c.subcomponents["header"].(common.Header).Height

		c.subcomponents["table"] = table
	}

	for name, cmp := range c.subcomponents {
		c.subcomponents[name], cmd = cmp.Update(msg)
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

func (c Lobby) View() string {
	layout := lipgloss.JoinVertical(lipgloss.Center,
		c.subcomponents["header"].View(),
		c.subcomponents["table"].View(),
	)

	return layout
}
