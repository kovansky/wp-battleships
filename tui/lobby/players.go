package lobby

import (
	"context"
	"github.com/76creates/stickers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/rs/zerolog"
	"unicode"
)

type Players struct {
	log zerolog.Logger

	theme tui.Theme

	Width  int
	Height int

	selected string

	table *stickers.Table
}

func CreatePlayers(ctx context.Context, theme tui.Theme) Players {
	table := stickers.NewTable(0, 0, []string{
		"Nickname",
		"Description",
		"Wins",
	})

	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	return Players{
		log:   log,
		theme: theme,
		table: table,
	}
}

func (c Players) Init() tea.Cmd {
	// Types
	var (
		s string
		i int
	)
	c.table, _ = c.table.SetTypes([]any{s, s, i}...)
	c.table.SetRatio([]int{3, 6, 1}).SetMinWidth([]int{5, 10, 3})

	c.table.SetStyles(map[stickers.TableStyleKey]lipgloss.Style{
		stickers.TableCellCursorStyleKey: lipgloss.NewStyle().
			Background(c.theme.TextPrimary.GetForeground()).
			Foreground(lipgloss.Color("#383838")),
		stickers.TableRowsCursorStyleKey: lipgloss.NewStyle(),
	})

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
		case "up":
			c.table.CursorUp()
		case "enter", " ":
			c.selected = c.table.GetCursorValue()
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
		var (
			playersTable [][]any
			err          error
		)

		for _, player := range msg.Players {
			playersTable = append(playersTable, []any{player.Name(), player.Description(), player.Wins()})
		}

		c.table, err = c.table.AddRows(playersTable)
		if err != nil {
			c.log.Error().Err(err).Msg("failed to add rows to table")
		}
	}

	return c, nil
}

func (c Players) View() string {
	return c.table.Render()
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
