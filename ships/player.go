package ships

import battleships "github.com/kovansky/wp-battleships"

var _ battleships.Player = (*Player)(nil)

type Player struct {
	name        string
	description string
	wins        int
}

func NewPlayer(name, description string) *Player {
	return &Player{name: name, description: description}
}

func NewPlayerWithWins(name, description string, wins int) *Player {
	return &Player{name: name, description: description, wins: wins}
}

func (p Player) Name() string {
	return p.name
}

func (p Player) Description() string {
	return p.description
}

func (p Player) Wins() int {
	return p.wins
}
