package battleships

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ContextKeyLog string = "battleships_logger"
)

var (
	Version      string
	ServerClient Client
	GameInstance Game

	ProgramMessage func(msg tea.Msg)

	Themes GameThemes
)
