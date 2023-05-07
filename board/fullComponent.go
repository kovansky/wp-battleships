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

type Full struct {
	themes themes

	friendly     Single
	opponent     Single
	playersInfo  string
	displayError string

	flexbox     *stickers.FlexBox
	asciiRender *figlet4go.AsciiRender

	targetInput textinput.Model

	battleships.Game
}

func InitFull(game battleships.Game, themeFriendly, themeEnemy, themeGlobal Theme, playersInfo string) Full {
	friendly := InitSingle(themeFriendly, game.Board())
	opponent := InitSingle(themeEnemy, game.OpponentBoard())
	flexbox := stickers.NewFlexBox(0, 0)
	asciiRender := figlet4go.NewAsciiRender()

	targetInput := textinput.New()
	targetInput.Placeholder = "Where to shoot, captain?"
	targetInput.CharLimit = 3
	targetInput.Width = 25

	return Full{
		themes:      themes{themeFriendly, themeEnemy, themeGlobal},
		friendly:    friendly,
		opponent:    opponent,
		playersInfo: playersInfo,
		flexbox:     flexbox,
		asciiRender: asciiRender,
		targetInput: targetInput,
		Game:        game,
	}
}

func (c Full) Init() tea.Cmd {
	c.friendly.Init()
	c.opponent.Init()

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

func (c Full) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			shotState, err := battleships.ServerClient.Fire(c.Game, field)
			if err != nil {
				c.displayError = "Error firing: " + err.Error()
				break
			}
			var fieldState battleships.FieldState
			switch shotState {
			case battleships.ShotMiss:
				fieldState = battleships.FieldStateMiss
			case battleships.ShotHit:
				fieldState = battleships.FieldStateHit
			case battleships.ShotSunk:
				fieldState = battleships.FieldStateSunk
			}

			board := c.OpponentBoard()
			if board == nil {
				board = make(map[string]battleships.FieldState)
			}

			board[field] = fieldState

			c.SetOpponentBoard(board)
			c.opponent.SetBoard(c.OpponentBoard())
		}
	case tea.WindowSizeMsg:
		c.flexbox.SetWidth(msg.Width)
		c.flexbox.SetHeight(msg.Height)
	case battleships.GameUpdateMsg:
		if c.GameStatus().ShouldFire {
			cmds = append(cmds, c.targetInput.Focus())
		}
	case battleships.PlayersUpdateMsg:
		c.playersInfo = msg.PlayersInfo
	}

	c.friendly, cmd = c.friendly.Update(msg)
	cmds = append(cmds, cmd)
	c.opponent, cmd = c.opponent.Update(msg)
	cmds = append(cmds, cmd)
	c.targetInput, cmd = c.targetInput.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c Full) View() string {
	friendlyRender, _ := c.asciiRender.Render("Friendly")
	enemyRender, _ := c.asciiRender.Render("Enemy")
	gameInfoRender, _ := c.asciiRender.Render("Game Info")

	friendlyState := c.themes.global.TextPrimary
	enemyState := c.themes.global.TextPrimary

	gameInfo := c.playersInfo

	if c.GameStatus().ShouldFire {
		friendlyState = c.themes.global.TextSecondary

		gameInfo += "\n\nYour turn!\n\t" + c.themes.global.TextSecondary.Render(strconv.Itoa(c.GameStatus().Timer)) + " seconds left to fire"
	} else {
		enemyState = c.themes.global.TextSecondary
	}

	c.flexbox.Row(0).Cell(0).SetContent(friendlyState.Render(friendlyRender))
	c.flexbox.Row(0).Cell(1).SetContent(enemyState.Render(enemyRender))
	c.flexbox.Row(0).Cell(2).SetContent(c.themes.global.TextPrimary.Render(gameInfoRender))

	c.flexbox.Row(1).Cell(0).SetContent(c.friendly.View())
	c.flexbox.Row(1).Cell(1).SetContent(c.opponent.View())
	c.flexbox.Row(1).Cell(2).SetContent(gameInfo)

	if c.GameStatus().ShouldFire {
		c.flexbox.Row(2).Cell(0).SetContent("\n\n\n" + c.targetInput.View() + "\n" + c.themes.global.TextSecondary.Render(c.displayError))
	}

	return c.flexbox.Render()
}

func fieldWithinBoard(field string) bool {
	if len(field) > 3 || len(field) < 2 {
		return false
	}

	if field[0] < 'A' || field[0] > 'J' {
		return false
	}

	if field[1] < '1' || field[1] > '9' {
		return false
	} else if len(field) == 3 && (field[1] != '1' || field[2] != '0') {
		return false
	}

	return true
}
