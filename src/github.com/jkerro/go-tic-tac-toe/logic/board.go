package logic

import "github.com/jkerro/go-tic-tac-toe/repository"

type BoardElement string

const (
	FREE BoardElement = "free"
	X    BoardElement = "x"
	O    BoardElement = "o"
)

type Board struct {
	Elements []BoardElement
}

func (b *Board) Row(index int) []BoardElement {
	return b.Elements[index*3 : index*3+3]
}

func (b *Board) Column(index int) []BoardElement {
	return []BoardElement{b.Elements[index], b.Elements[index+3], b.Elements[index+6]}
}

func (b *Board) Diagonals() ([]BoardElement, []BoardElement) {
	return []BoardElement{b.Elements[0], b.Elements[4], b.Elements[8]},
		[]BoardElement{b.Elements[2], b.Elements[4], b.Elements[6]}
}

func (b *Board) Matrix() [][]BoardElement {
	return [][]BoardElement{
		b.Row(0),
		b.Row(1),
		b.Row(2),
	}
}

func GetBoard(moves []repository.Move) *Board {
	b := NewBoard()
	for _, move := range moves {
		x := move.Col
		y := move.Row
		boardIndex := y*3 + x
		b.Elements[boardIndex] = GetTurnElement(move.Inx)
	}
	return b
}

func GetTurnElement(moveIndex int) BoardElement {
	if moveIndex%2 == 0 {
		return X
	} else {
		return O
	}
}

func NewBoard() *Board {
	b := &Board{make([]BoardElement, 9)}
	for i := 0; i < 9; i++ {
		b.Elements[i] = FREE
	}
	return b
}
