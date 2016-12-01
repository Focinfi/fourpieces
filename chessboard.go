package fourpieces

import (
	"encoding/json"
	"fmt"
	"strings"

	"strconv"

	"github.com/boltdb/bolt"
)

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

	player1 *Player
	player2 *Player
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

	game.player1 = newPlayer(PlayerA, game)
	game.player2 = newPlayer(PlayerB, game)
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
	nextX := step.ChessPiece.X + step.Direction.X
	nextY := step.ChessPiece.Y + step.Direction.Y
	// real piece
	if step.player.Type != game.board[step.ChessPiece.X][step.ChessPiece.Y] {
		return errStepInvalidPiece
	}

	if !inRange(nextX, nextY) {
		return errStepOutOfRange
	}

	// if free
	if game.board[nextX][nextY] != 0 {
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
	db, err := bolt.Open(fmt.Sprintf("/Users/Frank/work/go/src/github.com/Focinfi/fourpieces/player%s.games.data", player.Type), 0600, nil)
	if err != nil {
		return fmt.Errorf("bolt: %s", err.Error())
	}
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
				score++
			} else if game.winner == game.rivalOfPlayer(player).Type {
				score--
			}

			err = bucket.Put(b, []byte(strconv.Itoa(score)))

			if err != nil {
				return fmt.Errorf("save game: %s", err.Error())
			}

		}
		return nil
	})

	return nil
}
