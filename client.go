package battleships

type Client interface {
	InitGame(data GamePost) (Game, error)
}

type GamePost struct {
	Coords     []string `json:"coords,omitempty"`
	Desc       string   `json:"desc,omitempty"`
	Nick       string   `json:"nick,omitempty"`
	TargetNick string   `json:"target_nick,omitempty"`
	Wpbot      bool     `json:"wpbot,omitempty"`
}
