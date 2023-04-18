package board

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slices"
	"strings"
)

type Component struct {
	board []string
}

func InitialBoard(board []string) *Component {
	return &Component{board: board}
}

func (c Component) Init() tea.Cmd {
	return nil
}

func (c Component) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {
		case "ctrl+c", "q":
			return c, tea.Quit
		}
	}
	return c, nil
}

func (c Component) View() string {
	boardBuilder := strings.Builder{}

	cols, rows := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	for _, row := range rows {
		if row == "10" {
			boardBuilder.WriteString(row + " ")
		} else {
			boardBuilder.WriteString(row + "  ")
		}

		for _, col := range cols {
			if slices.Contains(c.board, fmt.Sprintf("%s%s", col, row)) {
				boardBuilder.WriteString("X  ")
			} else {
				boardBuilder.WriteString("   ")
			}
		}
		boardBuilder.WriteByte('\n')
	}

	boardBuilder.WriteString("   ")
	for _, col := range cols {
		boardBuilder.WriteString(col + "  ")
	}

	return boardBuilder.String()
}
