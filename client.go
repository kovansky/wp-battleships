package battleships

type Client interface {
	InitGame(data GamePost) (Game, error)
	Abandon(game Game) error

	UpdateBoard(game Game) error

	GameStatus(game Game) error
	GameDesc(game Game) error

	Refresh(game Game) error

	PlayerStats(nick string) (PlayerStats, error)
	ListPlayers() ([]Player, error)

	Fire(game Game, field string) (ShotState, error)
}

type Status string

const (
	StatusWaitingWPBot   Status = "waiting_wpbot"
	StatusWaiting               = "waiting"
	StatusGameInProgress        = "game_in_progress"
	StatusEnded                 = "ended"
	StatusWin                   = "win"
	StatusLose                  = "lose"
)

type ShotState string

const (
	ShotMiss ShotState = "miss"
	ShotHit            = "hit"
	ShotSunk           = "sunk"
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

type FireRes struct {
	Result ShotState `json:"result"`
}

type PlayerStats struct {
	Nick string `json:"nick"`

	Games  int `json:"games"`
	Wins   int `json:"wins"`
	Points int `json:"points"`

	Rank int `json:"rank"`
}
