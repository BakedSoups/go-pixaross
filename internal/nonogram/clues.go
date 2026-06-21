package nonogram

func RowClues(solution [][]bool) [][]int {
	clues := make([][]int, len(solution))
	for y := range solution {
		clues[y] = lineClues(solution[y])
	}
	return clues
}

func ColumnClues(solution [][]bool) [][]int {
	if len(solution) == 0 {
		return nil
	}

	height := len(solution)
	width := len(solution[0])
	clues := make([][]int, width)
	for x := 0; x < width; x++ {
		line := make([]bool, height)
		for y := 0; y < height; y++ {
			line[y] = solution[y][x]
		}
		clues[x] = lineClues(line)
	}
	return clues
}

func lineClues(line []bool) []int {
	var clues []int
	run := 0
	for _, filled := range line {
		if filled {
			run++
			continue
		}
		if run > 0 {
			clues = append(clues, run)
			run = 0
		}
	}
	if run > 0 {
		clues = append(clues, run)
	}
	if len(clues) == 0 {
		return []int{0}
	}
	return clues
}
