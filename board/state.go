package board

type State string

const (
	FieldEmpty     State = "empty"
	FieldShip            = "ship"
	FieldMiss            = "miss"
	FieldProtected       = "protected"
	FieldCorner          = "corner"
)

func (s State) Priority() int {
	switch s {
	case FieldEmpty:
		return 0
	case FieldShip:
		return 2
	case FieldProtected, FieldCorner:
		return 3
	default:
		return 1
	}
}
