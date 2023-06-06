package ranking

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

type Ranking struct {
	log zerolog.Logger

	theme battleships.Theme

	subcomponents  map[string]tea.Model
	initialPlayers []battleships.Player

	asciiRender *figlet4go.AsciiRender
}

func Create(ctx context.Context, theme battleships.Theme, players []battleships.Player) Ranking {
	asciiRender := figlet4go.NewAsciiRender()
	header := common.CreateHeader("Battleships", theme, asciiRender)
	table := CreateTable(ctx, theme, players)

	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	return Ranking{
		log:         log,
		theme:       theme,
		asciiRender: asciiRender,
		subcomponents: map[string]tea.Model{
			"header": header,
			"table":  table,
		},
		initialPlayers: players,
	}
}

func (c Ranking) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, cmp := range c.subcomponents {
		cmds = append(cmds, cmp.Init())
	}

	return tea.Batch(cmds...)
}

func (c Ranking) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return c, tea.Quit
		case "L", "l":
			return c, func() tea.Msg {
				return tui.ApplicationStageChangeMsg{
					From:  tui.StageLogin,
					Stage: tui.StageLobby,
				}
			}
		}
	case tea.WindowSizeMsg:
		for name, cmp := range c.subcomponents {
			c.subcomponents[name], cmd = cmp.Update(msg)
			cmds = append(cmds, cmd)
		}

		table := c.subcomponents["table"].(Table)

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

func (c Ranking) View() string {
	layout := lipgloss.JoinVertical(lipgloss.Center,
		c.subcomponents["header"].View(),
		c.subcomponents["table"].View(),
	)

	return layout
}
