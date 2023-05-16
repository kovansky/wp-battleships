package battleships

type Statistics struct {
	shots int
	hits  int
	sunk  int
}

func NewStatistics() *Statistics {
	return &Statistics{
		shots: 0,
		hits:  0,
		sunk:  0,
	}
}

func (s *Statistics) Shots() int {
	return s.shots
}

func (s *Statistics) Hits() int {
	return s.hits
}

func (s *Statistics) Sunk() int {
	return s.sunk
}

func (s *Statistics) SetShots(count int) {
	s.shots = count
}

func (s *Statistics) SetHits(count int) {
	s.hits = count
}

func (s *Statistics) SetSunk(count int) {
	s.sunk = count
}

func (s *Statistics) IncrementShots() {
	s.shots++
}

func (s *Statistics) IncrementHits() {
	s.hits++
}

func (s *Statistics) IncrementSunk() {
	s.sunk++
}
