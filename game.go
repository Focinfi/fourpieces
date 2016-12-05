package fourpieces

import (
	"errors"
	"fmt"
	"time"

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
	if game.playerA.turnNum != game.currentTurn ||
		game.playerB.turnNum != game.currentTurn {
		game.err = errPlayerNotSync
		game.setOver()
	}

	defer func() {
		game.currentTurn++
		game.playerA.turnNum++
		game.playerB.turnNum++
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
	aPiecesCnt := len(game.playerA.pieces)
	bPiecesCnt := len(game.playerB.pieces)

	// 3-1 3 win
	if aPiecesCnt >= 3 && aPiecesCnt-bPiecesCnt >= 2 {
		game.winner = game.playerA.Type
		game.setOver()
	} else if bPiecesCnt >= 3 && bPiecesCnt-aPiecesCnt >= 2 {
		game.winner = game.playerB.Type
		game.setOver()
	}

	// draw
	if !game.isOver() {
		// currentTurn > 20000
		if game.currentTurn > 20000 {
			game.winner = 0
			game.setOver()
		} else if (aPiecesCnt == 2 && bPiecesCnt == 1 ||
			aPiecesCnt == 1 && bPiecesCnt == 2) && game.currentTurn > 10000 {
			// 1-2 or 2-1
			game.winner = 0
			game.setOver()
		} else if (aPiecesCnt == 2 && bPiecesCnt == 2) && game.currentTurn > 10000 {
			// 2-2
			game.winner = 0
			game.setOver()
		}
	}
}

func (game *chessBoard) nextStep(t PlayerType) error {
	var step *Step
	switch t {
	case game.playerA.Type:
		step = game.playerA.nextStep()
	case game.playerB.Type:
		step = game.playerB.nextStep()
	default:
		return errStepInvalidPiece
	}

	if step == nil {
		game.winner = rivalOfPlayerType(t)
		game.setOver()
		return nil
	}

	err := game.applyStep(step)
	if err != nil {
		fmt.Println(err)
		return err
	}

	game.checkOver()

	fmt.Printf("player1: %d, player2: %d\n", len(game.playerA.pieces), len(game.playerB.pieces))
	fmt.Println(game)
	time.Sleep(time.Second)
	return nil
}

func (game *chessBoard) applyStep(step *Step) error {
	game.board[step.basePiece.X][step.basePiece.Y] = 0
	game.board[step.MoveTo.X][step.MoveTo.Y] = step.player.Type
	step.basePiece.X = step.MoveTo.X
	step.basePiece.Y = step.MoveTo.Y

	eatedPieces, err := game.eatedPiece(step.MoveTo, step.player)
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
		fmt.Printf("eated Player[%v]: piece(%d, %d)\n", rival.Type, toRemove.X, toRemove.Y)
		for i, piece := range rival.pieces {
			if piece.X == toRemove.X && piece.Y == toRemove.Y {
				game.board[toRemove.X][toRemove.Y] = 0
				rival.pieces = append(rival.pieces[:i], rival.pieces[i+1:]...)
				// fmt.Printf("pieces in player: %v\n", rival.pieces)
				rival.steps[len(rival.steps)-1].score -= 10
				fmt.Printf("reduce score: %#v\n", rival.steps[len(rival.steps)-1].MoveTo)
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
		rival = game.playerB
	} else if player.Type == PlayerB {
		rival = game.playerA
	}

	return
}

func (game *chessBoard) eatedPiece(piece *ChessPiece, player *Player) (eated []ChessPiece, err error) {
	rivalType := game.rivalOfPlayer(player).Type
	return eatPieces(piece.X, piece.Y, player.Type, rivalType, game.board)
}

func eatPieces(pieceX, pieceY int, playerType, rivalType PlayerType, board [][]PlayerType) (eated []ChessPiece, err error) {
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
		// for _, piece := range game.playerA.pieces {

		// 	fmt.Printf("(%d, %d)\n", piece.X, piece.Y)
		// }

		// cmd := ""
		// for cmd != "y" {
		// 	fmt.Printf("\nNext turn?(y)")
		// 	fmt.Scanln(&cmd)
		// }
	}
	if game.err != nil {
		log.Fatal("game: " + game.err.Error())

	}

	if err := game.saveToFS(game.playerA); err != nil {
		log.Fatal(err)
	}

	if err := game.saveToFS(game.playerB); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("winner: %v, turn: %d\n", game.winner, game.currentTurn)
}
