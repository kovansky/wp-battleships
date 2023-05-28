package wait

import (
	"context"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/tui/common"
	"github.com/mbndr/figlet4go"
	"github.com/rs/zerolog"
)

type Wait struct {
	ctx context.Context
	log zerolog.Logger

	theme battleships.Theme

	subcomponents map[string]tea.Model

	asciiRender *figlet4go.AsciiRender
}

func Create(ctx context.Context, theme battleships.Theme) Wait {
	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	asciiRender := figlet4go.NewAsciiRender()

	header := common.CreateHeader("Battleships", theme, asciiRender)
	spinnerComponent := CreateSpinner()
	spinnerComponent.Spinner = spinner.Points
	spinnerComponent.Style = battleships.Themes.Global.TextSecondary()

	return Wait{
		ctx:   ctx,
		log:   log,
		theme: theme,
		subcomponents: map[string]tea.Model{
			"header":  header,
			"spinner": spinnerComponent,
		},
		asciiRender: asciiRender,
	}
}

func (c Wait) Init() tea.Cmd {
	return c.subcomponents["spinner"].(Spinner).Model.Tick
}

func (c Wait) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return c, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		c.subcomponents["spinner"], cmd = c.subcomponents["spinner"].(Spinner).Update(msg)
		return c, cmd
	}
	
	for name, cmp := range c.subcomponents {
		c.subcomponents[name], cmd = cmp.Update(msg)
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

func (c Wait) View() string {
	block := lipgloss.JoinVertical(lipgloss.Center,
		c.theme.TextPrimary().
			Copy().
			Italic(true).
			Render("Waiting for challenge...\n"),
		c.subcomponents["spinner"].View(),
	)

	layout := lipgloss.JoinVertical(lipgloss.Center,
		c.subcomponents["header"].View(),
		lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			Padding(1, 3).
			Render(block),
	)

	return lipgloss.JoinHorizontal(lipgloss.Center, layout)
}
