package fourpieces

// HEIGHT for chess board max height
const HEIGHT = 3

// WEIGHT for chess board max weight
const WEIGHT = 3

func inRange(x, y int) bool {
	return x >= 0 && x <= HEIGHT && y >= 0 && y <= WEIGHT
}
