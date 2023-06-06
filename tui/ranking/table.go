package ranking

import (
	"context"
	"github.com/76creates/stickers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/rs/zerolog"
	"unicode"
)

type Table struct {
	log zerolog.Logger

	theme battleships.Theme

	Width  int
	Height int

	focused             int
	selected            string
	initialRankingTable []battleships.Player

	filterString string

	table *stickers.Table
}

func CreateTable(ctx context.Context, theme battleships.Theme, initialRankingTable []battleships.Player) Table {
	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	table, err := initializeTable(theme, initialRankingTable...)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize table")
	}

	return Table{
		log:   log,
		theme: theme,
		table: table,
	}
}

func (c Table) Init() tea.Cmd {
	return nil
}

func (c Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return c, tea.Quit
		case "down":
			c.table.CursorDown()
			_, c.focused = c.table.GetCursorLocation()
		case "up":
			c.table.CursorUp()
			_, c.focused = c.table.GetCursorLocation()
		case "backspace":
			c.filterWithStr(msg.String())
			_, c.filterString = c.table.GetFilter()
		default:
			if len(msg.String()) == 1 {
				r := msg.Runes[0]
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					c.filterWithStr(msg.String())
					_, c.filterString = c.table.GetFilter()
				}
			}
		}
	case tea.WindowSizeMsg:
		c.table.
			SetWidth(c.Width).
			SetHeight(c.Height)
	}

	return c, nil
}

func (c Table) View() string {
	return c.table.Render()
}

func initializeTable(theme battleships.Theme, players ...battleships.Player) (*stickers.Table, error) {
	table := stickers.NewTable(0, 0, []string{
		"Position",
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
	table, _ = table.SetTypes([]any{i, s, i, i, i}...)
	table.SetRatio([]int{1, 6, 2, 1, 1}).SetMinWidth([]int{2, 10, 4, 3, 3})

	table.SetStyles(map[stickers.TableStyleKey]lipgloss.Style{
		stickers.TableCellCursorStyleKey: lipgloss.NewStyle().
			Background(theme.TextPrimary().GetForeground()).
			Foreground(lipgloss.Color("#383838")),
		stickers.TableRowsCursorStyleKey: lipgloss.NewStyle(),
	})

	var (
		RankingTableTable [][]any
		err               error
	)

	for id, player := range players {
		RankingTableTable = append(RankingTableTable, []any{id + 1, player.Name(), player.Points(), player.Wins(), player.Games()})
	}

	table, err = table.AddRows(RankingTableTable)

	return table, err
}

func (c Table) filterWithStr(key string) {
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
