package battleships

type Player interface {
	Name() string
	Description() string
	SetWins(wins int)
	Wins() int
	SetPoints(points int)
	Points() int
	SetGames(points int)
	Games() int
	SetRank(points int)
	Rank() int
}
