package fourpieces

// HEIGHT for chess board max height
const HEIGHT = 3

// WEIGHT for chess board max weight
const WEIGHT = 3

func inRange(x, y int) bool {
	return x >= 0 && x <= HEIGHT && y >= 0 && y <= WEIGHT
}

// eatPieces returns which pieces should be eated.
// pieceX, pieceY is the chess piece of the player with type playerType moved,
//
func eatPieces(pieceX, pieceY int, playerType PlayerType, board [][]PlayerType) (eated []ChessPiece, err error) {
	if board[pieceX][pieceY] == 0 {
		err = errUnknownPiece
		return
	}

	rivalType := rivalOfPlayerType(playerType)
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
