package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ApplicationStageChangeMsg struct {
	From  Stage
	Stage Stage
	Model tea.Model
}
