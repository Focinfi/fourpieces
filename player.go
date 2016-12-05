package fourpieces

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/boltdb/bolt"
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

func rivalOfPlayerType(t PlayerType) PlayerType {
	if t == PlayerA {
		return PlayerB
	} else if t == PlayerB {
		return PlayerA
	}

	return 0
}

// Player is for a player
type Player struct {
	Type PlayerType

	game    *chessBoard
	pieces  []*ChessPiece
	steps   []*Step
	exDB    *bolt.DB
	turnNum int
}

func newPlayer(t PlayerType, game *chessBoard) *Player {
	db, err := bolt.Open(dataPath(t), 0600, nil)
	if err != nil {
		log.Fatalf("newPlayer: can not load experirence db, err: %v", err)
	}

	return &Player{
		Type:   t,
		game:   game,
		pieces: newChessPieces(t),
		exDB:   db,
	}
}

func (player *Player) nextStep() *Step {
	stepOtps := player.otpSteps()
	if len(stepOtps) <= 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	stepIdx := rand.Intn(len(stepOtps))
	fmt.Printf("Player[% 2v], opt[%d], move piece %v\n", player.Type, len(stepOtps), stepOtps[stepIdx].MoveTo)
	step := stepOtps[stepIdx]
	player.steps = append(player.steps, step)
	// fmt.Printf("%v\n%v => %v\n", step.Board, step.ChessPiece, step.Direction)
	// time.Sleep(time.Second)
	return step
}

func (player *Player) availibleSteps() (steps []*Step) {
	for _, piece := range player.pieces {
		for _, direction := range stepDirections {
			step := newStep(player, piece, direction)
			// check position availability
			if err := player.game.checkStepPosition(*step); err == nil {
				step.Board = moveOneStepOnBoard(step.Board, step)
				step.setScore()
				steps = append(steps, step)
			}
		}
	}
	return
}

func (player *Player) otpSteps() (steps []*Step) {
	avlSteps := player.availibleSteps()

	// advance 1 step
	for _, step := range avlSteps {
		n := step.checkNextStep()
		step.score += n * 10
		fmt.Printf("otpSteps: basePiece: %v, MoveTo: %v, score: %v\n", step.basePiece, step.MoveTo, step.score)
	}

	// collect the steps with top score
	for _, step := range avlSteps {
		if len(steps) == 0 {
			steps = append(steps, step)
			continue
		}

		topScore := steps[len(steps)-1].score
		if step.score == topScore {
			steps = append(steps, step)
		} else if step.score > topScore {
			steps = []*Step{step}
		}
	}

	return
}
