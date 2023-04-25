package board

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type Component struct {
	theme Theme
	ships map[string]interface{}
}

func InitComponent(theme Theme, ships ...string) Component {
	shipsMap := make(map[string]interface{}, len(ships))

	for _, ship := range ships {
		shipsMap[strings.ToUpper(ship)] = nil
	}

	return Component{theme: theme, ships: shipsMap}
}

func (c Component) Init() tea.Cmd {
	return nil
}

func (c Component) Update(msg tea.Msg) (Component, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return c, tea.Quit
		}
	}

	return c, nil
}

func (c Component) View() string {
	cols, rows := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}, []string{"10", "9", "8", "7", "6", "5", "4", "3", "2", "1"}
	const sep = " "
	builder := strings.Builder{}

	for _, rowLabel := range rows {
		label := rowLabel

		if len(rowLabel) == 1 {
			label = sep + rowLabel
		}

		builder.WriteString(c.theme.Rows.Render(label))
		builder.WriteString(strings.Repeat(sep, 1))

		for _, colLabel := range cols {
			field := colLabel + rowLabel

			if _, contains := c.ships[field]; contains {
				builder.WriteString(c.theme.RenderShip())
				continue
			}

			builder.WriteString(strings.Repeat(sep, 3))
		}

		builder.WriteByte('\n')
	}

	builder.WriteString(strings.Repeat(sep, 3))
	for _, colLabel := range cols {
		builder.WriteString(sep + c.theme.Cols.Render(colLabel) + sep)
	}

	return builder.String()
}
