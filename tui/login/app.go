package login

import (
	"context"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/kovansky/wp-battleships/tui/common"
	"github.com/kovansky/wp-battleships/tui/lobby"
	"github.com/mbndr/figlet4go"
	"github.com/rs/zerolog"
)

type Login struct {
	ctx context.Context
	log zerolog.Logger

	theme battleships.Theme

	subcomponents map[string]tea.Model
	inputs        []textinput.Model

	focusIndex int

	asciiRender *figlet4go.AsciiRender
}

func Create(ctx context.Context, theme battleships.Theme) Login {
	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	asciiRender := figlet4go.NewAsciiRender()

	inputComponents := make([]textinput.Model, 2)
	inputs := []struct {
		name      string
		charLimit int
	}{
		{"nick", 30},
		{"description", 150},
	}
	for i, data := range inputs {
		input := textinput.New()
		input.CharLimit = data.charLimit
		input.Width = 50

		if i == 0 {
			input.Focus()
		} else {
			input.Blur()
		}

		inputComponents[i] = input
	}

	header := common.CreateHeader("Battleships", theme, asciiRender)
	submitButton := common.CreateButton("Submit", theme)

	return Login{
		ctx:    ctx,
		log:    log,
		theme:  theme,
		inputs: inputComponents,
		subcomponents: map[string]tea.Model{
			"header": header,
			"submit": submitButton,
		},
		asciiRender: asciiRender,
	}
}

func (c Login) Init() tea.Cmd {
	return textinput.Blink
}

func (c Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tmp  textinput.Model
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return c, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Is user trying to submit the whole thing?
			if s == "enter" && c.focusIndex == len(c.inputs) {
				return c, c.submit()
			}

			// Is user trying to move up/down?
			if s == "up" || s == "shift+tab" {
				c.focusIndex--
			} else {
				c.focusIndex++
			}

			// Is user trying to move out of bounds?
			if c.focusIndex < 0 {
				c.focusIndex = len(c.inputs)
			} else if c.focusIndex > len(c.inputs) {
				c.focusIndex = 0
			}

			for i := 0; i < len(c.inputs); i++ {
				if i == c.focusIndex {
					cmds = append(cmds, c.inputs[i].Focus())
					continue
				}

				c.inputs[i].Blur()
			}

			if c.focusIndex == len(c.inputs) {
				c.subcomponents["submit"] = c.subcomponents["submit"].(common.Button).Focus()
			} else {
				c.subcomponents["submit"] = c.subcomponents["submit"].(common.Button).Blur()
			}

			return c, tea.Batch(cmds...)
		}
	}

	for name, cmp := range c.subcomponents {
		c.subcomponents[name], cmd = cmp.Update(msg)
		cmds = append(cmds, cmd)
	}
	for name, cmp := range c.inputs {
		tmp, cmd = cmp.Update(msg)
		c.inputs[name] = tmp
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

func (c Login) View() string {
	block := lipgloss.JoinVertical(lipgloss.Center,
		c.theme.TextPrimary().
			Copy().
			Italic(true).
			Render("Enter your nickname and description.\n"+
				"You may leave them blank - they will be randomly selected for you c:\n"),
		// Nick
		c.theme.TextSecondary().Render("Nickname:"),
		c.inputs[0].View(),
		// Description
		c.theme.TextSecondary().Render("Description:"),
		c.inputs[1].View(),
		// Submit button
		"",
		c.subcomponents["submit"].View(),
	)

	layout := lipgloss.JoinVertical(lipgloss.Center,
		c.subcomponents["header"].View(),
		lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			Padding(1, 3).
			Render(block),
	)

	return layout
}

func (c Login) submit() tea.Cmd {
	battleships.PlayerData.Nick = c.inputs[0].Value()
	battleships.PlayerData.Description = c.inputs[1].Value()

	players, err := battleships.ServerClient.ListPlayers()
	if err != nil {
		c.log.Fatal().Err(err).Msg("failed to list players")
	}

	lobbyApp := lobby.Create(c.ctx, battleships.Themes.Global, players)

	return func() tea.Msg {
		return tui.ApplicationStageChangeMsg{
			From:  tui.StageLogin,
			Stage: tui.StageLobby,
			Model: lobbyApp,
		}
	}
}
