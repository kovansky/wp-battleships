package tui

import (
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
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

func (b Brush) SetChar(char byte) Brush {
	b.char = char
	return b
}

func (b Brush) Style() lipgloss.Style {
	return b.style
}

func (b Brush) SetStyle(style lipgloss.Style) Brush {
	b.style = style
	return b
}

type Theme struct {
	Rows          lipgloss.Style
	Cols          lipgloss.Style
	TextPrimary   lipgloss.Style
	TextSecondary lipgloss.Style
	Border        Brush
	Ship          Brush
	Hit           Brush
	Sunk          Brush
	Miss          Brush
}

func NewTheme() Theme {
	return Theme{}
}

func (t Theme) SetRows(style lipgloss.Style) Theme {
	t.Rows = style
	return t
}

func (t Theme) SetCols(style lipgloss.Style) Theme {
	t.Cols = style
	return t
}

func (t Theme) SetTextPrimary(style lipgloss.Style) Theme {
	t.TextPrimary = style
	return t
}

func (t Theme) SetTextSecondary(style lipgloss.Style) Theme {
	t.TextSecondary = style
	return t
}

func (t Theme) SetBorder(brush Brush) Theme {
	brush = brush.SetStyle(brush.Style().PaddingRight(1))
	t.Border = brush
	return t
}

func (t Theme) SetShip(brush Brush) Theme {
	brush = brush.SetStyle(brush.Style().PaddingRight(1))
	t.Ship = brush
	return t
}

func (t Theme) SetHit(brush Brush) Theme {
	brush = brush.SetStyle(brush.Style().PaddingRight(1))
	t.Hit = brush
	return t
}

func (t Theme) SetSunk(brush Brush) Theme {
	brush = brush.SetStyle(brush.Style().PaddingRight(1))
	t.Sunk = brush
	return t
}

func (t Theme) SetMiss(brush Brush) Theme {
	brush = brush.SetStyle(brush.Style().PaddingRight(1))
	t.Miss = brush
	return t
}

func (t Theme) RenderBorder() string {
	return t.Border.Style().Render(string(t.Border.Char()))
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

func (t Theme) RenderShip() string {
	return t.Ship.Style().Render(string(t.Ship.Char()))
}

func (t Theme) RenderHit() string {
	return t.Hit.Style().Render(string(t.Hit.Char()))
}

func (t Theme) RenderSunk() string {
	return t.Sunk.Style().Render(string(t.Sunk.Char()))
}

func (t Theme) RenderMiss() string {
	return t.Miss.Style().Render(string(t.Miss.Char()))
}
