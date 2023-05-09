package lobby

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/mbndr/figlet4go"
)

type Header struct {
	Text  string
	theme tui.Theme

	blockStyle lipgloss.Style

	width  int
	Height int

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
		c.blockStyle = c.theme.TextPrimary.Copy().
			Width(c.width).
			Align(lipgloss.Center).
			PaddingBottom(2)

		ascii, _ := c.asciiRender.Render(c.Text)

		render := c.blockStyle.Render(ascii)

		c.Height = lipgloss.Height(render)
	}

	return c, nil
}

func (c Header) View() string {
	ascii, _ := c.asciiRender.Render(c.Text)

	render := c.blockStyle.Render(ascii)

	return render
}
