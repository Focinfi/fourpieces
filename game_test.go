package fourpieces

import (
	"testing"

	"github.com/boltdb/bolt"
)

func Test(t *testing.T) {
	game := newFourPieces()
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
	game := newFourPieces()
	player2 := game.rivalOfPlayer(game.playerA)
	if player2 != game.playerB {
		t.Error("can not get eh rival player")
	}

	player2.pieces = nil

	if game.playerB.pieces != nil {
		t.Error("can not change the pieces")
	}
}

func TestSaveToFS(t *testing.T) {
	db, err := bolt.Open(dataPath(PlayerB), 0600, nil)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		steps := tx.Bucket([]byte("steps"))
		if steps == nil {
			t.Error("empty database")
		}

		score := steps.Get([]byte(`{"MoveTo":{"X":2,"Y":2},"Board":[[1,0,0,-1],[1,0,0,0],[1,0,-1,0],[0,1,0,-1]]}`))
		t.Log(string(score))
		return nil
	})
}
