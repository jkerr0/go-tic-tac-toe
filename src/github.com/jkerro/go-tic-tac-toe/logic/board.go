package logic

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
	return b.Elements[index : index+3]
}

func (b *Board) Matrix() [][]BoardElement {
	return [][]BoardElement{
		b.Row(0),
		b.Row(1),
		b.Row(2),
	}
}

func NewBoard() *Board {
	b := &Board{make([]BoardElement, 9)}
	for i := 0; i < 9; i++ {
		b.Elements[i] = FREE
	}
	return b
}
