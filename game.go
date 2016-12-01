package fourpieces

import (
	"errors"
	"fmt"

	"log"
)

var errPlayerNotSync = errors.New("game: player(s) out of contorl")
var errUnknownPiece = errors.New("rule: nuknown piece")

func (game chessBoard) isOver() bool {
	return game.over
}

func (game *chessBoard) setOver() {
	game.over = true
}

func (game *chessBoard) nextTurn() {
	// check state
	if game.player1.turnNum != game.currentTurn ||
		game.player2.turnNum != game.currentTurn {
		game.err = errPlayerNotSync
		game.setOver()
	}

	defer func() {
		game.currentTurn++
		game.player1.turnNum++
		game.player2.turnNum++
	}()

	// PLAYER1 first
	if err := game.nextStep(PlayerA); err != nil {
		game.setOver()
		return
	}

	if !game.isOver() {
		// PLAYER2 second
		if err := game.nextStep(PlayerB); err != nil {
			game.setOver()
			return
		}
	}

	return
}

func (game *chessBoard) checkOver() {
	if len(game.player1.pieces)-len(game.player2.pieces) >= 2 {
		game.winner = game.player1.Type
		game.setOver()
	} else if len(game.player2.pieces)-len(game.player1.pieces) >= 2 {
		game.winner = game.player2.Type
		game.setOver()
	}

	// draw
	if len(game.player1.pieces) == 2 && len(game.player2.pieces) == 1 ||
		len(game.player1.pieces) == 1 && len(game.player2.pieces) == 2 {
		game.winner = 0
		game.setOver()
	}
}

func (game *chessBoard) nextStep(t PlayerType) error {
	var step *Step
	switch t {
	case game.player1.Type:
		step = game.player1.nextStep()
	case game.player2.Type:
		step = game.player2.nextStep()
	default:
		return errStepInvalidPiece
	}

	if step == nil {
		game.winner = game.rivalOfPlayer(step.player).Type
		game.setOver()
		return nil
	}

	err := game.applyStep(step)
	if err != nil {
		return err
	}

	game.checkOver()

	fmt.Printf("player1: %d, player2: %d\n", len(game.player1.pieces), len(game.player2.pieces))
	// time.Sleep(time.Second)
	return nil
}

func (game *chessBoard) applyStep(step *Step) error {
	game.board[step.ChessPiece.X][step.ChessPiece.Y] = 0
	step.ChessPiece.moveStep(step.Direction)
	game.board[step.ChessPiece.X][step.ChessPiece.Y] = step.player.Type

	eatedPieces, err := game.eatedPiece(step.ChessPiece, step.player)
	if err != nil {
		return err
	}

	if err := game.removeRivalPieces(eatedPieces, step.player); err != nil {
		return err
	}

	return nil
}

func (game *chessBoard) removeRivalPieces(pieces []ChessPiece, player *Player) error {
	rival := game.rivalOfPlayer(player)
	for _, toRemove := range pieces {
		fmt.Printf("eated Player[%d]: piece(%d, %d)\n", rival.Type, toRemove.X, toRemove.Y)
		for i, piece := range rival.pieces {
			if piece.X == toRemove.X && piece.Y == toRemove.Y {
				game.board[toRemove.X][toRemove.Y] = 0
				rival.pieces = append(rival.pieces[:i], rival.pieces[i+1:]...)
				// fmt.Printf("pieces in player: %v\n", rival.pieces)
			}
		}
	}

	if len(rival.pieces) == 0 {
		game.winner = player.Type
		game.setOver()
	}

	return nil
}

func (game *chessBoard) rivalOfPlayer(player *Player) (rival *Player) {
	if player.Type == PlayerA {
		rival = game.player2
	} else if player.Type == PlayerB {
		rival = game.player1
	}

	return
}

func (game *chessBoard) eatedPiece(piece *ChessPiece, player *Player) (eated []ChessPiece, err error) {
	pieceX := piece.X
	pieceY := piece.Y
	board := game.board
	playerType := player.Type
	rivalType := game.rivalOfPlayer(player).Type

	if board[pieceX][pieceY] == 0 {
		err = errUnknownPiece
		return
	}

	// x line
	eatedXLineRivalIdx := -1
	xLine := board[pieceX]
	if xLine[0] == 0 && xLine[1] == rivalType && xLine[2] == playerType && xLine[3] == playerType {
		eatedXLineRivalIdx = 1
	} else if xLine[0] == 0 && xLine[1] == playerType && xLine[2] == playerType && xLine[3] == rivalType {
		eatedXLineRivalIdx = 3
	} else if xLine[0] == rivalType && xLine[1] == playerType && xLine[2] == playerType && xLine[3] == 0 {
		eatedXLineRivalIdx = 0
	} else if xLine[0] == playerType && xLine[1] == playerType && xLine[2] == rivalType && xLine[3] == 0 {
		eatedXLineRivalIdx = 2
	}
	if eatedXLineRivalIdx != -1 {
		eated = append(eated, ChessPiece{X: pieceX, Y: eatedXLineRivalIdx})
	}

	// y line
	eatedYLineRivalIdx := -1
	yLine := func(i int) PlayerType {
		return board[i][pieceY]
	}
	if yLine(0) == 0 && yLine(1) == rivalType && yLine(2) == playerType && yLine(3) == playerType {
		eatedYLineRivalIdx = 1
	} else if yLine(0) == 0 && yLine(1) == playerType && yLine(2) == playerType && yLine(3) == rivalType {
		eatedYLineRivalIdx = 3
	} else if yLine(0) == rivalType && yLine(1) == playerType && yLine(2) == playerType && yLine(3) == 0 {
		eatedYLineRivalIdx = 0
	} else if yLine(0) == playerType && yLine(1) == playerType && yLine(2) == rivalType && yLine(3) == 0 {
		eatedYLineRivalIdx = 2
	}
	if eatedYLineRivalIdx != -1 {
		eated = append(eated, ChessPiece{X: eatedYLineRivalIdx, Y: pieceY})
	}

	return
}

// Play play for testing
func Play() {
	game := newChessBoard()
	for !game.isOver() {
		game.nextTurn()
	}
	if game.err != nil {
		log.Fatal("game: " + game.err.Error())

	}

	if err := game.saveToFS(game.player1); err != nil {
		log.Fatal(err)
	}

	if err := game.saveToFS(game.player2); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("winner: %v, turn: %d", game.winner, game.currentTurn)
}
