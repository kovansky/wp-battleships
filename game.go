package battleships

type Game interface {
	SetPlayer(player Player)
	Player() Player
	Key() string
}
