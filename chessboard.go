package fourpieces

import (
	"encoding/json"
	"fmt"
	"strings"

	"strconv"

	"os"
	"path"

	"github.com/boltdb/bolt"
)

var appDir = path.Join(os.Getenv("GOPATH"), "src", "github.com", "Focinfi", "fourpieces")

// ChessPiece as a chess piece
type ChessPiece struct{ X, Y int }

func (piece *ChessPiece) moveStep(direction StepDirection) {
	piece.X += direction.X
	piece.Y += direction.Y
}

// newChessPieces allocates and return a new []*ChessPiece,
// which contains four pieces.
// if player is 1, pieces are (0, 0), (1, 0), (2, 0), (3, 0),
// otherwise, pieces are (0, 3), (1, 3), (2, 3), (3, 3).
func newChessPieces(player PlayerType) (pieces []*ChessPiece) {
	y := 0
	if player != PlayerA {
		y = 3
	}

	for x := 0; x <= 3; x++ {
		pieces = append(pieces, &ChessPiece{x, y})
	}

	return
}

type chessBoard struct {
	id          int
	currentTurn int
	over        bool
	err         error

	playerA *Player
	playerB *Player
	board   [][]PlayerType

	winner PlayerType
}

func newChessBoard() chessBoard {
	game := &chessBoard{
		id: genChessBoardID(),
		board: [][]PlayerType{
			{PlayerA, 0, 0, PlayerB},
			{PlayerA, 0, 0, PlayerB},
			{PlayerA, 0, 0, PlayerB},
			{PlayerA, 0, 0, PlayerB},
		},
	}

	game.playerA = newPlayer(PlayerA, game)
	game.playerB = newPlayer(PlayerB, game)
	return *game
}

func genChessBoardID() int {
	return 1
}

func (game chessBoard) boardSnapshot() [][]PlayerType {
	board := make([][]PlayerType, 4)
	for x := 0; x <= HEIGHT; x++ {
		board[x] = make([]PlayerType, 4)
		for y := 0; y <= WEIGHT; y++ {
			board[x][y] = game.board[x][y]
		}
	}
	return board
}

func (game chessBoard) checkStepPosition(step Step) error {
	// real piece
	if step.player.Type != game.board[step.basePiece.X][step.basePiece.Y] {
		return errStepInvalidPiece
	}

	if !inRange(step.MoveTo.X, step.MoveTo.Y) {
		return errStepOutOfRange
	}

	// if free
	if game.board[step.MoveTo.X][step.MoveTo.Y] != 0 {
		return errStepNoFree
	}

	return nil
}

func (game chessBoard) String() string {
	lines := []string{fmt.Sprintf("> T%d\n", game.currentTurn)}
	for _, xLine := range game.board {
		lines = append(lines, fmt.Sprintf("% 2v\n", xLine))
	}

	lines = append(lines, "\n")

	return strings.Join(lines, "")
}

type board [][]PlayerType

func (b board) piece(x, y int) PlayerType {
	if y > len(b) {
		return -1
	}

	// yLine := b[x]
	// if x > len(yLine)

	return 0
}

func newBoard() [][]PlayerType {
	return [][]PlayerType{
		{PlayerA, 0, 0, PlayerB},
		{PlayerA, 0, 0, PlayerB},
		{PlayerA, 0, 0, PlayerB},
		{PlayerA, 0, 0, PlayerB},
	}
}

func (game chessBoard) saveToFS(player *Player) error {
	db := player.exDB
	defer db.Close()

	db.Batch(func(tx *bolt.Tx) error {
		for _, step := range player.steps {
			b, err := json.Marshal(step)
			if err != nil {
				return fmt.Errorf("save game: %s", game.err.Error())
			}
			bucket, err := tx.CreateBucketIfNotExists([]byte("steps"))
			if err != nil {
				return fmt.Errorf("save game: %s", err.Error())
			}

			score := 0
			scoreBytes := bucket.Get(b)
			if scoreBytes != nil {
				score, err = strconv.Atoi(string(scoreBytes))
				if err != nil {
					continue
				}
			}

			if game.winner == step.player.Type {
				score += 2
			} else if game.winner == game.rivalOfPlayer(player).Type {
				score--
			} else {
				score++
			}

			err = bucket.Put(b, []byte(strconv.Itoa(score+step.score)))

			if err != nil {
				return fmt.Errorf("save game: %s", err.Error())
			}

		}
		return nil
	})

	return nil
}

func dataPath(t PlayerType) string {
	return path.Join(appDir, fmt.Sprintf("player%s.games.data", t))
}

func moveOneStepOnBoard(board [][]PlayerType, step *Step) [][]PlayerType {
	fmt.Println(board, step.basePiece, step.MoveTo)
	board[step.basePiece.X][step.basePiece.Y] = 0
	board[step.MoveTo.X][step.MoveTo.Y] = step.player.Type
	return board
}
