package lobby

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/mbndr/figlet4go"
)

type Header struct {
	Text        string
	theme       tui.Theme
	width       int
	asciiRender *figlet4go.AsciiRender
}

func CreateHeader(text string, theme tui.Theme, asciiRender *figlet4go.AsciiRender) Header {
	return Header{
		Text:        text,
		theme:       theme,
		asciiRender: asciiRender,
	}
}

func (c Header) Init() tea.Cmd {
	return nil
}

func (c Header) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
	}

	return c, nil
}

func (c Header) View() string {
	block := c.theme.TextPrimary.Copy().
		Width(c.width).
		Align(lipgloss.Center).
		PaddingBottom(2)

	ascii, _ := c.asciiRender.Render(c.Text)

	return block.Render(ascii)
}
