package ships

import battleships "github.com/kovansky/wp-battleships"

var _ battleships.Player = (*Player)(nil)

type Player struct {
	name        string
	description string
	wins        int
	games       int
	points      int
}

func NewPlayer(name, description string) *Player {
	return &Player{name: name, description: description}
}

func NewPlayerFromStats(stats battleships.PlayerStats) *Player {
	return &Player{
		name:   stats.Nick,
		wins:   stats.Wins,
		games:  stats.Games,
		points: stats.Points,
	}
}

func (p Player) Name() string {
	return p.name
}

func (p Player) Description() string {
	return p.description
}

func (p Player) SetWins(wins int) {
	p.wins = wins
}

func (p Player) Wins() int {
	return p.wins
}

func (p Player) SetPoints(points int) {
	p.points = points
}

func (p Player) Points() int {
	return p.points
}

func (p Player) SetGames(games int) {
	p.games = games
}

func (p Player) Games() int {
	return p.games
}
