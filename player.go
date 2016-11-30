package fourpieces

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// PlayerType for role of the two players
type PlayerType int

const (
	// PlayerA go first
	PlayerA = 1
	// PlayerB go second
	PlayerB = -1
)

func (t PlayerType) String() (s string) {
	switch t {
	case PlayerA:
		s = "A"
	case PlayerB:
		s = "B"
	default:
		s = "-"
	}

	return
}

func isRival(t1, t2 PlayerType) bool {
	return math.Abs(float64(t1-t2)) == float64(math.Abs(float64(PlayerA)-float64(PlayerB)))
}

// Player is for a player
type Player struct {
	Type PlayerType

	game    *chessBoard
	pieces  []*ChessPiece
	turnNum int
}

func newPlayer(t PlayerType, game *chessBoard) *Player {
	return &Player{
		Type:   t,
		game:   game,
		pieces: newChessPieces(t),
	}
}

func (player *Player) nextStep() *Step {
	stepOtps := player.availableSteps()
	if len(stepOtps) > 0 {
		rand.Seed(time.Now().Unix())
		stepIdx := rand.Intn(len(stepOtps))
		fmt.Printf("Player[% 2d], opt[%d], move piece[%d]\n", player.Type, len(stepOtps), stepOtps[stepIdx].chessPiece.x)
		return stepOtps[stepIdx]
	}

	return nil
}

func (player *Player) availableSteps() (steps []*Step) {
	for _, piece := range player.pieces {
		for _, direction := range stepDirections {
			step := Step{
				player:     player,
				chessPiece: piece,
				direction:  direction,
			}

			// check position availability
			if err := player.game.checkStepPosition(step); err == nil {
				steps = append(steps, &step)
			}
		}
	}

	return
}
