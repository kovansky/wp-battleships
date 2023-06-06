package tui

type Stage string

const (
	StageLogin Stage = "login"
	StageSetup       = "setup"
	StageWait        = "wait"
	StageLobby       = "lobby"
	StageGame        = "game"

	StageRanking = "ranking"
)
