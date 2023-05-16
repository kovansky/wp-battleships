package battleships

const (
	ContextKeyLog string = "battleships_logger"
)

var (
	Version      string
	ServerClient Client
	GameInstance Game
)
