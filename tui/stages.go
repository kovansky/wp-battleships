package tui

type Stage string

const (
	StageLogin Stage = "login"
	StageWait        = "wait"
	StageLobby       = "lobby"
	StageGame        = "game"
)
