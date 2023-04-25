package board

import (
	"github.com/76creates/stickers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mbndr/figlet4go"
)

type themes struct {
	friendly Theme
	enemy    Theme
	global   Theme
}

type FullComponent struct {
	themes themes

	friendly Component
	enemy    Component

	flexbox     *stickers.FlexBox
	asciiRender *figlet4go.AsciiRender
}

func InitFullComponent(themeFriendly, themeEnemy, themeGlobal Theme, shipsFriendly, shipsEnemy []string) FullComponent {
	friendly := InitComponent(themeFriendly, shipsFriendly...)
	enemy := InitComponent(themeEnemy, shipsEnemy...)
	flexbox := stickers.NewFlexBox(0, 0)
	asciiRender := figlet4go.NewAsciiRender()

	return FullComponent{themes{
		themeFriendly,
		themeEnemy,
		themeGlobal,
	}, friendly, enemy, flexbox, asciiRender}
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
	}

	c.flexbox.AddRows(rows)

	return nil
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
		}
	case tea.WindowSizeMsg:
		c.flexbox.SetWidth(msg.Width)
		c.flexbox.SetHeight(msg.Height)
	}

	c.friendly, cmd = c.friendly.Update(msg)
	cmds = append(cmds, cmd)
	c.enemy, cmd = c.enemy.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c FullComponent) View() string {
	friendlyRender, _ := c.asciiRender.Render("Friendly")
	enemyRender, _ := c.asciiRender.Render("Enemy")
	gameInfoRender, _ := c.asciiRender.Render("Game Info")

	c.flexbox.Row(0).Cell(0).SetContent(c.themes.global.Text.Render(friendlyRender))
	c.flexbox.Row(0).Cell(1).SetContent(c.themes.global.Text.Render(enemyRender))
	c.flexbox.Row(0).Cell(2).SetContent(c.themes.global.Text.Render(gameInfoRender))

	c.flexbox.Row(1).Cell(0).SetContent(c.friendly.View())
	c.flexbox.Row(1).Cell(1).SetContent(c.enemy.View())

	return c.flexbox.Render()
}
