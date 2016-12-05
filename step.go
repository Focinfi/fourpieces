package fourpieces

import (
	"errors"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

import "encoding/json"

import "strconv"

// StepDirection for one step direction
type StepDirection struct {
	X, Y int
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
	player    *Player
	score     int
	basePiece *ChessPiece

	MoveTo *ChessPiece
	Board  [][]PlayerType
}

var errStepInvalidPiece = errors.New("step: unknown piece")
var errStepOutOfRange = errors.New("step: out of chess board range")
var errStepNoFree = errors.New("step: no free room")

func newStep(player *Player, basePiece *ChessPiece, direction StepDirection) *Step {
	step := &Step{
		player:    player,
		basePiece: basePiece,
		MoveTo:    &ChessPiece{X: basePiece.X + direction.X, Y: basePiece.Y + direction.Y},
		Board:     player.game.boardSnapshot(),
	}
	return step
}

func (step *Step) setScore() {
	step.player.exDB.View(func(tx *bolt.Tx) error {
		steps := tx.Bucket([]byte("steps"))
		if steps != nil {
			scoreBytes := steps.Get(step.toJSONBytes())
			score, err := strconv.Atoi(string(scoreBytes))
			if scoreBytes == nil {
				fmt.Println(string(step.toJSONBytes()))
			}

			if scoreBytes != nil && err != nil {
				log.Printf("step: score saved in database is not int, err: %v\n", err)
			}
			step.score = score
		}
		return nil
	})

}

func (step Step) toJSONBytes() []byte {
	b, err := json.Marshal(step)
	if err != nil {
		log.Fatalf("step: can not marshal into JSON, err: %v\n", err)
	}

	return b
}

func (step Step) checkNextStep() int {
	board := step.Board
	board[step.basePiece.X][step.basePiece.Y] = 0
	board[step.MoveTo.X][step.MoveTo.Y] = step.player.Type

	eated, _ := eatPieces(step.MoveTo.X,
		step.MoveTo.Y,
		step.player.Type,
		step.player.game.rivalOfPlayer(step.player).Type,
		board)
	return len(eated)
}
