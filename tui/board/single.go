package board

import (
	tea "github.com/charmbracelet/bubbletea"
	battleships "github.com/kovansky/wp-battleships"
	"strings"
)

type Single struct {
	theme  battleships.Theme
	fields map[string]battleships.FieldState
}

func InitSingle(theme battleships.Theme, board map[string]battleships.FieldState) Single {
	return Single{theme: theme, fields: board}
}

func (c *Single) Init() tea.Cmd {
	return nil
}

func (c *Single) Update(msg tea.Msg) (Single, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return *c, tea.Quit
		}
	}

	return *c, nil
}

func (c *Single) SetBoard(board map[string]battleships.FieldState) {
	c.fields = board
}

func (c *Single) View() string {
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
				builder.WriteString(c.theme.RenderField(state))
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
