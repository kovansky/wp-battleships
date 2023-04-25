package battleships

type Client interface {
	InitGame(data GamePost) (Game, error)
	UpdateBoard(game Game) error
	GameStatus(game Game) error
	GameDesc(game Game) error
}

type Status string

const (
	StatusWaitingWPBot   Status = "waiting_wpbot"
	StatusWaiting        Status = "waiting"
	StatusGameInProgress Status = "game_in_progress"
	StatusEnded          Status = "ended"
	StatusWin            Status = "win"
	StatusLose           Status = "lose"
)

type GamePost struct {
	Coords     []string `json:"coords,omitempty"`
	Desc       string   `json:"desc,omitempty"`
	Nick       string   `json:"nick,omitempty"`
	TargetNick string   `json:"target_nick,omitempty"`
	Wpbot      bool     `json:"wpbot,omitempty"`
}

type GameGet struct {
	Nick string `json:"nick"`
	Desc string `json:"desc"`

	Opponent string   `json:"opponent"`
	OppDesc  string   `json:"opp_desc"`
	OppShots []string `json:"opp_shots"`

	GameStatus     Status `json:"game_status"`
	LastGameStatus Status `json:"last_game_status"`

	ShouldFire bool `json:"should_fire"`
	Timer      int  `json:"timer"`
}

type BoardGet struct {
	Board []string `json:"board"`
}
