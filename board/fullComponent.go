package board

import (
	"github.com/76creates/stickers"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/mbndr/figlet4go"
	"strconv"
	"strings"
)

type themes struct {
	friendly Theme
	enemy    Theme
	global   Theme
}

type FullComponent struct {
	themes themes

	friendly     Component
	enemy        Component
	playersInfo  string
	displayError string

	flexbox     *stickers.FlexBox
	asciiRender *figlet4go.AsciiRender

	targetInput textinput.Model

	battleships.Game
}

func InitFullComponent(game battleships.Game, themeFriendly, themeEnemy, themeGlobal Theme, playersInfo string, shipsFriendly, shipsEnemy []string) FullComponent {
	friendly := InitComponent(themeFriendly, shipsFriendly...)
	enemy := InitComponent(themeEnemy, shipsEnemy...)
	flexbox := stickers.NewFlexBox(0, 0)
	asciiRender := figlet4go.NewAsciiRender()

	targetInput := textinput.New()
	targetInput.Placeholder = "Where to shoot, captain?"
	targetInput.CharLimit = 3
	targetInput.Width = 25

	return FullComponent{
		themes:      themes{themeFriendly, themeEnemy, themeGlobal},
		friendly:    friendly,
		enemy:       enemy,
		playersInfo: playersInfo,
		flexbox:     flexbox,
		asciiRender: asciiRender,
		targetInput: targetInput,
		Game:        game,
	}
}

func (c FullComponent) Init() tea.Cmd {
	c.friendly.Init()
	c.enemy.Init()

	rows := []*stickers.FlexBoxRow{
		// Headers row
		c.flexbox.NewRow().AddCells(
			[]*stickers.FlexBoxCell{
				stickers.NewFlexBoxCell(1, 1),
				stickers.NewFlexBoxCell(1, 1),
				stickers.NewFlexBoxCell(1, 1),
			},
		),
		// Boards row
		c.flexbox.NewRow().AddCells(
			[]*stickers.FlexBoxCell{
				stickers.NewFlexBoxCell(1, 1),
				stickers.NewFlexBoxCell(1, 1),
				stickers.NewFlexBoxCell(1, 1),
			},
		),
		c.flexbox.NewRow().AddCells(
			[]*stickers.FlexBoxCell{
				stickers.NewFlexBoxCell(1, 1),
			}),
	}
	c.flexbox.AddRows(rows)

	return textinput.Blink
}

func (c FullComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return c, tea.Quit
		case "enter":
			field := strings.ToUpper(c.targetInput.Value())
			if !fieldWithinBoard(field) {
				c.displayError = "Field outside of board"
				break
			}
			c.displayError = ""
			c.targetInput.SetValue("")

			_, err := battleships.ServerClient.Fire(c.Game, field)
			if err != nil {
				c.displayError = "Error firing: " + err.Error()
				break
			}

			c.Update(battleships.GameUpdateMsg{})
		}
	case tea.WindowSizeMsg:
		c.flexbox.SetWidth(msg.Width)
		c.flexbox.SetHeight(msg.Height)
	case battleships.GameUpdateMsg:
		if c.Game.GameStatus().ShouldFire {
			cmds = append(cmds, c.targetInput.Focus())
		}
	case battleships.PlayersUpdateMsg:
		c.playersInfo = msg.PlayersInfo
	}

	c.friendly, cmd = c.friendly.Update(msg)
	cmds = append(cmds, cmd)
	c.enemy, cmd = c.enemy.Update(msg)
	cmds = append(cmds, cmd)
	c.targetInput, cmd = c.targetInput.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c FullComponent) View() string {
	friendlyRender, _ := c.asciiRender.Render("Friendly")
	enemyRender, _ := c.asciiRender.Render("Enemy")
	gameInfoRender, _ := c.asciiRender.Render("Game Info")

	friendlyState := c.themes.global.TextPrimary
	enemyState := c.themes.global.TextPrimary

	gameInfo := c.playersInfo

	if c.Game.GameStatus().ShouldFire {
		friendlyState = c.themes.global.TextSecondary

		gameInfo += "\n\nYour turn!\n\t" + c.themes.global.TextSecondary.Render(strconv.Itoa(c.Game.GameStatus().Timer)) + " seconds left to fire"
	} else {
		enemyState = c.themes.global.TextSecondary
	}

	c.flexbox.Row(0).Cell(0).SetContent(friendlyState.Render(friendlyRender))
	c.flexbox.Row(0).Cell(1).SetContent(enemyState.Render(enemyRender))
	c.flexbox.Row(0).Cell(2).SetContent(c.themes.global.TextPrimary.Render(gameInfoRender))

	c.flexbox.Row(1).Cell(0).SetContent(c.friendly.View())
	c.flexbox.Row(1).Cell(1).SetContent(c.enemy.View())
	c.flexbox.Row(1).Cell(2).SetContent(gameInfo)

	if c.Game.GameStatus().ShouldFire {
		c.flexbox.Row(2).Cell(0).SetContent("\n\n\n" + c.targetInput.View() + "\n" + c.themes.global.TextSecondary.Render(c.displayError))
	}

	return c.flexbox.Render()
}

func fieldWithinBoard(field string) bool {
	if len(field) != 2 {
		return false
	}

	if field[0] < 'A' || field[0] > 'J' {
		return false
	}

	if field[1] < '0' || field[1] > '9' {
		return false
	}

	return true
}
