package lobby

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/mbndr/figlet4go"
)

type Lobby struct {
	theme tui.Theme

	subcomponents map[string]tea.Model

	asciiRender *figlet4go.AsciiRender
}

func Create(theme tui.Theme) Lobby {
	asciiRender := figlet4go.NewAsciiRender()
	header := CreateHeader("Battleships", theme, asciiRender)

	return Lobby{
		theme:       theme,
		asciiRender: asciiRender,
		subcomponents: map[string]tea.Model{
			"header": header,
		},
	}
}

func (c Lobby) Init() tea.Cmd {
	return nil
}

func (c Lobby) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return c, tea.Quit
		}
	}

	for name, cmp := range c.subcomponents {
		c.subcomponents[name], cmd = cmp.Update(msg)
		cmds = append(cmds, cmd)
	}

	return c, nil
}

func (c Lobby) View() string {
	return c.subcomponents["header"].View()
}
