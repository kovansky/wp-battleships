package tui

import (
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/parts"
)

type Brush struct {
	char  byte
	style lipgloss.Style
}

func NewBrush() Brush {
	return Brush{}
}

func (b Brush) Char() byte {
	return b.char
}

func (b Brush) SetChar(char byte) battleships.Brush {
	b.char = char
	return b
}

func (b Brush) Style() lipgloss.Style {
	return b.style
}

func (b Brush) SetStyle(style lipgloss.Style) battleships.Brush {
	b.style = style
	return b
}

type Theme struct {
	rows          lipgloss.Style
	cols          lipgloss.Style
	textPrimary   lipgloss.Style
	textSecondary lipgloss.Style
	border        battleships.Brush
	ship          battleships.Brush
	hit           battleships.Brush
	sunk          battleships.Brush
	miss          battleships.Brush
}

func NewTheme() battleships.Theme {
	return Theme{}
}

func (t Theme) Rows() lipgloss.Style {
	return t.rows
}

func (t Theme) SetRows(style lipgloss.Style) battleships.Theme {
	t.rows = style
	return t
}

func (t Theme) Cols() lipgloss.Style {
	return t.cols
}

func (t Theme) SetCols(style lipgloss.Style) battleships.Theme {
	t.cols = style
	return t
}

func (t Theme) TextPrimary() lipgloss.Style {
	return t.textPrimary
}

func (t Theme) SetTextPrimary(style lipgloss.Style) battleships.Theme {
	t.textPrimary = style
	return t
}

func (t Theme) TextSecondary() lipgloss.Style {
	return t.textSecondary
}

func (t Theme) SetTextSecondary(style lipgloss.Style) battleships.Theme {
	t.textSecondary = style
	return t
}

func (t Theme) Border() battleships.Brush {
	return t.border
}

func (t Theme) SetBorder(brush battleships.Brush) battleships.Theme {
	brush = brush.SetStyle(brush.Style().PaddingRight(1))
	t.border = brush
	return t
}

func (t Theme) Ship() battleships.Brush {
	return t.ship

}
func (t Theme) SetShip(brush battleships.Brush) battleships.Theme {
	brush = brush.SetStyle(brush.Style().PaddingRight(1))
	t.ship = brush
	return t
}

func (t Theme) Hit() battleships.Brush {
	return t.hit

}
func (t Theme) SetHit(brush battleships.Brush) battleships.Theme {
	brush = brush.SetStyle(brush.Style().PaddingRight(1))
	t.hit = brush
	return t
}

func (t Theme) Sunk() battleships.Brush {
	return t.sunk

}
func (t Theme) SetSunk(brush battleships.Brush) battleships.Theme {
	brush = brush.SetStyle(brush.Style().PaddingRight(1))
	t.sunk = brush
	return t
}

func (t Theme) Miss() battleships.Brush {
	return t.miss
}
func (t Theme) SetMiss(brush battleships.Brush) battleships.Theme {
	brush = brush.SetStyle(brush.Style().PaddingRight(1))
	t.miss = brush
	return t
}

func (t Theme) RenderBorder() string {
	return t.border.Style().Render(string(t.border.Char()))
}

func (t Theme) RenderField(state battleships.FieldState) string {
	switch state {
	case battleships.FieldStateHit:
		return t.RenderHit()
	case battleships.FieldStateMiss:
		return t.RenderMiss()
	case battleships.FieldStateShip:
		return t.RenderShip()
	case battleships.FieldStateSunk:
		return t.RenderSunk()
	default:
		return "  "
	}
}

func (t Theme) NewRenderField(state parts.State) string {
	switch state {
	case parts.FieldMiss, parts.FieldProtected, parts.FieldCorner:
		return t.RenderMiss()
	case parts.FieldShip:
		return t.RenderShip()
	default:
		return "  "
	}
}

func (t Theme) RenderShip() string {
	return t.ship.Style().Render(string(t.ship.Char()))
}

func (t Theme) RenderHit() string {
	return t.hit.Style().Render(string(t.hit.Char()))
}

func (t Theme) RenderSunk() string {
	return t.sunk.Style().Render(string(t.sunk.Char()))
}

func (t Theme) RenderMiss() string {
	return t.miss.Style().Render(string(t.miss.Char()))
}
