package battleships

import "github.com/charmbracelet/lipgloss"

type GameThemes struct {
	Player Theme
	Enemy  Theme
	Global Theme
}

type Brush interface {
	Char() byte
	SetChar(char byte) Brush
	Style() lipgloss.Style
	SetStyle(style lipgloss.Style) Brush
}

type Theme interface {
	Rows() lipgloss.Style
	SetRows(style lipgloss.Style) Theme
	Cols() lipgloss.Style
	SetCols(style lipgloss.Style) Theme

	TextPrimary() lipgloss.Style
	SetTextPrimary(style lipgloss.Style) Theme
	TextSecondary() lipgloss.Style
	SetTextSecondary(style lipgloss.Style) Theme

	Border() Brush
	SetBorder(brush Brush) Theme
	Ship() Brush
	SetShip(brush Brush) Theme
	Hit() Brush
	SetHit(brush Brush) Theme
	Sunk() Brush
	SetSunk(brush Brush) Theme
	Miss() Brush
	SetMiss(brush Brush) Theme

	RenderBorder() string
	RenderField(state FieldState) string

	RenderShip() string
	RenderHit() string
	RenderSunk() string
	RenderMiss() string
}
