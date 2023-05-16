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

	Themes   GameThemes
	Routines GameRoutines

	ProgramMessage func(msg tea.Msg)
)
