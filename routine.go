package battleships

type GameRoutines struct {
	Lobby Routine
	Game  Routine
	Wait  Routine
}

type Routine interface {
	Run()
	Quit()
}
