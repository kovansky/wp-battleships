package wait

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Spinner struct {
	spinner.Model
}

func CreateSpinner() Spinner {
	return Spinner{
		Model: spinner.New(),
	}
}

func (c Spinner) Init() tea.Cmd {
	return nil
}

func (c Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	c.Model, cmd = c.Model.Update(msg)

	return c, cmd
}

func (c Spinner) View() string {
	return c.Model.View()
}
