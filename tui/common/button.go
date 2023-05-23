package common

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
)

type Button struct {
	Text  string
	theme battleships.Theme

	Focused bool

	blockStyle lipgloss.Style
}

func CreateButton(text string, theme battleships.Theme) Button {
	return Button{
		Text:  text,
		theme: theme,
	}
}

func (c Button) Init() tea.Cmd {
	return nil
}

func (c Button) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c Button) View() string {
	renderer := c.theme.TextSecondary().Render
	if c.Focused == true {
		renderer = c.theme.TextPrimary().Copy().Underline(true).Render
	}

	render := fmt.Sprintf("[ %s ]", renderer(c.Text))

	return render
}

func (c Button) Focus() tea.Model {
	c.Focused = true

	return c
}

func (c Button) Blur() tea.Model {
	c.Focused = false

	return c
}
