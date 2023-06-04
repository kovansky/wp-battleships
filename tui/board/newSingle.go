package board

import (
	tea "github.com/charmbracelet/bubbletea"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/parts"
	"strings"
)

type NewSingle struct {
	theme  battleships.Theme
	fields map[string]parts.State
}

func InitNewSingle(theme battleships.Theme, board map[string]parts.State) NewSingle {
	return NewSingle{theme: theme, fields: board}
}

func (c *NewSingle) Init() tea.Cmd {
	return nil
}

func (c *NewSingle) Update(msg tea.Msg) (NewSingle, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return *c, tea.Quit
		}
	}

	return *c, nil
}

func (c *NewSingle) SetBoard(board map[string]parts.State) {
	c.fields = board
}

func (c *NewSingle) View() string {
	cols, rows := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}, []string{"10", "9", "8", "7", "6", "5", "4", "3", "2", "1"}
	const sep = " "
	builder := strings.Builder{}

	for _, rowLabel := range rows {
		label := rowLabel

		if len(rowLabel) == 1 {
			label = sep + rowLabel
		}

		builder.WriteString(c.theme.Rows().Render(label))
		builder.WriteString(strings.Repeat(sep, 1))

		for _, colLabel := range cols {
			field := colLabel + rowLabel

			if state, contains := c.fields[field]; contains {
				builder.WriteString(c.theme.NewRenderField(state))
				continue
			}

			builder.WriteString(strings.Repeat(sep, 2))
		}

		builder.WriteByte('\n')
	}

	builder.WriteString(strings.Repeat(sep, 3))
	for _, colLabel := range cols {
		builder.WriteString(c.theme.Cols().Render(colLabel) + sep)
	}

	return builder.String()
}
