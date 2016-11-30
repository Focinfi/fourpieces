package fourpieces

import (
	"testing"
)

func Test(t *testing.T) {
	game := newChessBoard()
	println(game.String())
	game.nextTurn()
	t.Log(game)
}

func TestEatedPiece(t *testing.T) {

}

func TestIsRival(t *testing.T) {
	if isRival(PlayerA, PlayerB) == true {
		t.Log("ok")
	} else {
		t.Error("fail")
	}

	if isRival(PlayerA, PlayerA) == false {
		t.Log("ok")
	} else {
		t.Error("fail")
	}

	if isRival(PlayerB, PlayerB) == false {
		t.Log("ok")
	} else {
		t.Error("fail")
	}
}

func TestPlayer(t *testing.T) {
	game := newChessBoard()
	player2 := game.rivalOfPlayer(game.player1)
	if player2 != game.player2 {
		t.Error("can not get eh rival player")
	}

	player2.pieces = nil

	if game.player2.pieces != nil {
		t.Error("can not change the pieces")
	}
}
