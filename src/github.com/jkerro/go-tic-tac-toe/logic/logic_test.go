package logic

import (
	"testing"
)

func TestWinOnDiagonal(t *testing.T) {
	board := NewBoard()
	board.Elements[0] = X
	board.Elements[4] = X
	board.Elements[8] = X
	win := CheckWin(board)
	if !win {
		t.Failed()
	}
}

func TestWinInColumn(t *testing.T) {
	board := NewBoard()
	board.Elements[0] = X
	board.Elements[3] = X
	board.Elements[6] = X
	win := CheckWin(board)
	if !win {
		t.FailNow()
	}
}

func TestWinInRow(t *testing.T) {
	board := NewBoard()
	board.Elements[0] = X
	board.Elements[1] = X
	board.Elements[2] = X
	win := CheckWin(board)
	if !win {
		t.FailNow()
	}
}

func TestStartWithX(t *testing.T) {
	turn := GetTurnElement(0)
	if turn != X {
		t.FailNow()
	}
}
