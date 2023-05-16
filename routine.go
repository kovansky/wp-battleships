package battleships

type GameRoutines struct {
	Lobby Routine
	Game  Routine
}

type Routine interface {
	Run()
	Quit()
}
