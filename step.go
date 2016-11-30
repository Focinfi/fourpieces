package fourpieces

import "errors"

// StepDirection for one step direction
type StepDirection struct {
	x, y int
}

var (
	// Up move one step in up side
	Up = StepDirection{-1, 0}
	// Down move one step in down side
	Down = StepDirection{1, 0}
	// Left move one step in left side
	Left = StepDirection{0, -1}
	// Right move one step in right side
	Right = StepDirection{0, 1}

	stepDirections = []StepDirection{Up, Down, Left, Right}
)

// Step for one step in the game
type Step struct {
	player     *Player
	chessPiece *ChessPiece
	direction  StepDirection
}

var errStepInvalidPiece = errors.New("step: unknown piece")
var errStepOutOfRange = errors.New("step: out of chess board range")
var errStepNoFree = errors.New("step: no free room")
