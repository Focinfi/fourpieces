package fourpieces

import (
	"fmt"
	"strings"
)

// ChessPiece as a chess piece
type ChessPiece struct{ x, y int }

func (piece *ChessPiece) moveStep(direction StepDirection) {
	piece.x += direction.x
	piece.y += direction.y
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

func (game chessBoard) checkStepPosition(step Step) error {
	nextX := step.chessPiece.x + step.direction.x
	nextY := step.chessPiece.y + step.direction.y
	// real piece
	if step.player.Type != game.board[step.chessPiece.x][step.chessPiece.y] {
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
