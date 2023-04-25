package board

import "github.com/charmbracelet/lipgloss"

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
	Rows   lipgloss.Style
	Cols   lipgloss.Style
	Text   lipgloss.Style
	Border Brush
	Ship   Brush
	Hit    Brush
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

func (t Theme) SetText(style lipgloss.Style) Theme {
	t.Text = style
	return t
}

func (t Theme) SetBorder(brush Brush) Theme {
	brush = brush.SetStyle(brush.Style().PaddingLeft(1).PaddingRight(1))
	t.Border = brush
	return t
}

func (t Theme) SetShip(brush Brush) Theme {
	brush = brush.SetStyle(brush.Style().PaddingLeft(1).PaddingRight(1))
	t.Ship = brush
	return t
}
func (t Theme) SetHit(brush Brush) Theme {
	brush = brush.SetStyle(brush.Style().PaddingLeft(1).PaddingRight(1))
	t.Hit = brush
	return t
}

func (t Theme) RenderBorder() string {
	return t.Border.Style().Render(string(t.Border.Char()))
}

func (t Theme) RenderShip() string {
	return t.Ship.Style().Render(string(t.Ship.Char()))
}

func (t Theme) RenderHit() string {
	return t.Hit.Style().Render(string(t.Hit.Char()))
}
