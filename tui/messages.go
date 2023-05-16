package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ApplicationStageChangeMsg struct {
	Stage Stage
	Board tea.Model
}
