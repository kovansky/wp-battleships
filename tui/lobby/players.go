package lobby

import (
	"context"
	"fmt"
	"github.com/76creates/stickers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/kovansky/wp-battleships/tui/board"
	"github.com/rs/zerolog"
	"unicode"
)

type Players struct {
	log zerolog.Logger

	theme tui.Theme

	Width  int
	Height int

	focused        int
	selected       string
	initialPlayers []battleships.Player

	table *stickers.Table
}

func CreatePlayers(ctx context.Context, theme tui.Theme, initialPlayers []battleships.Player) Players {
	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	table, err := initializeTable(theme, initialPlayers...)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize table")
	}

	return Players{
		log:   log,
		theme: theme,
		table: table,
	}
}

func (c Players) Init() tea.Cmd {
	return nil
}

func (c Players) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c", "ctrl+c":
			return c, tea.Quit
		case "down":
			c.table.CursorDown()
			_, c.focused = c.table.GetCursorLocation()
		case "up":
			c.table.CursorUp()
			_, c.focused = c.table.GetCursorLocation()
		case "enter", " ":
			c.selected = c.table.GetCursorValue()

			var (
				game battleships.Game
				err  error
			)
			if c.selected == "WP_Bot" {
				game, err = battleships.ServerClient.InitGame(battleships.GamePost{
					Wpbot: true,
				})
				if err != nil {
					c.log.Fatal().Err(err).Msg("Couldn't initialize game")
				}
			} else {
				game, err = battleships.ServerClient.InitGame(battleships.GamePost{
					TargetNick: c.selected,
				})
				if err != nil {
					c.log.Fatal().Err(err).Msg("Couldn't initialize game")
				}
			}

			battleships.GameInstance = game

			err = battleships.ServerClient.UpdateBoard(battleships.GameInstance)
			if err != nil {
				c.log.Fatal().Err(err).Msg("Couldn't update the game board")
			}
			err = battleships.ServerClient.GameStatus(battleships.GameInstance)
			if err != nil {
				c.log.Fatal().Err(err).Msg("Couldn't update the game status")
			}

			gameBoard := board.InitFull(game, c.theme, c.theme, c.theme, fmt.Sprintf(lipgloss.NewStyle().Italic(true).Render("Waiting for game...")))
			return c, func() tea.Msg {
				return tui.ApplicationStageChangeMsg{
					Stage: tui.StageGame,
					Board: gameBoard,
				}
			}
		case "backspace":
			c.filterWithStr(msg.String())
		default:
			if len(msg.String()) == 1 {
				r := msg.Runes[0]
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					c.filterWithStr(msg.String())
				}
			}
		}
	case tea.WindowSizeMsg:
		c.table.
			SetWidth(c.Width).
			SetHeight(c.Height)
	case battleships.PlayersListMsg:
		var err error
		c.table, err = initializeTable(c.theme, msg.Players...)
		if err != nil {
			c.log.Error().Err(err).Msg("failed to add rows to table")
		}
		c.table.
			SetWidth(c.Width).
			SetHeight(c.Height)

		// Move cursor to it's focused position
		for i := 0; i < c.focused; i++ {
			c.table.CursorDown()
		}
	}

	return c, nil
}

func (c Players) View() string {
	return c.table.Render()
}

func initializeTable(theme tui.Theme, players ...battleships.Player) (*stickers.Table, error) {
	table := stickers.NewTable(0, 0, []string{
		"Nickname",
		"Points",
		"Wins",
		"Games",
	})

	// Types
	var (
		s string
		i int
	)
	table, _ = table.SetTypes([]any{s, i, i, i}...)
	table.SetRatio([]int{6, 2, 1, 1}).SetMinWidth([]int{10, 4, 3, 3})

	table.SetStyles(map[stickers.TableStyleKey]lipgloss.Style{
		stickers.TableCellCursorStyleKey: lipgloss.NewStyle().
			Background(theme.TextPrimary.GetForeground()).
			Foreground(lipgloss.Color("#383838")),
		stickers.TableRowsCursorStyleKey: lipgloss.NewStyle(),
	})

	var (
		playersTable [][]any
		err          error
	)

	for _, player := range players {
		playersTable = append(playersTable, []any{player.Name(), player.Points(), player.Wins(), player.Games()})
	}

	table, err = table.AddRows(playersTable)

	return table, err
}

func (c Players) filterWithStr(key string) {
	i, s := c.table.GetFilter()
	x, _ := c.table.GetCursorLocation()
	if x != i && key != "backspace" {
		c.table.SetFilter(x, key)
		return
	}
	if key == "backspace" {
		if len(s) == 1 {
			c.table.UnsetFilter()
			return
		} else if len(s) > 1 {
			s = s[0 : len(s)-1]
		} else {
			return
		}
	} else {
		s = s + key
	}
	c.table.SetFilter(i, s)
}
