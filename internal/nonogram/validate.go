package nonogram

func IsSolved(board *Board, solution [][]bool) bool {
	if board == nil || len(solution) != board.Height {
		return false
	}
	for y := 0; y < board.Height; y++ {
		if len(solution[y]) != board.Width {
			return false
		}
		for x := 0; x < board.Width; x++ {
			if (board.Cells[y][x] == CellFilled) != solution[y][x] {
				return false
			}
		}
	}
	return true
}
