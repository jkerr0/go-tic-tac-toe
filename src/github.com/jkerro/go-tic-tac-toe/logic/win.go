package logic

func CheckWin(board *Board) bool {
	for i := 0; i < 2; i++ {
		row := board.Row(i)
		if checkElementsForWin(row) {
			return true
		}
	}

	for i := 0; i < 2; i++ {
		col := board.Column(i)
		if checkElementsForWin(col) {
			return true
		}
	}

	diag1, diag2 := board.Diagonals()
	if checkElementsForWin(diag1) {
		return true
	}
	if checkElementsForWin(diag2) {
		return true
	}
	return false
}

func checkElementsForWin(elements []BoardElement) bool {
	if elements[0] == FREE {
		return false
	}
	for i, curr := range elements[1:] {
		prev := elements[i]
		if curr != prev {
			return false
		}
	}
	return true
}
