package nonogram

import "testing"

func TestClues(t *testing.T) {
	solution := [][]bool{
		{false, true, true, false, true},
		{false, false, false, false, false},
		{true, true, true, true, true},
	}

	rows := RowClues(solution)
	assertClues(t, rows[0], []int{2, 1})
	assertClues(t, rows[1], []int{0})
	assertClues(t, rows[2], []int{5})

	cols := ColumnClues(solution)
	assertClues(t, cols[0], []int{1})
	assertClues(t, cols[1], []int{1, 1})
	assertClues(t, cols[3], []int{1})
}

func TestBoardApplyAndSolve(t *testing.T) {
	board := NewBoard(2, 2)
	solution := [][]bool{
		{true, false},
		{false, true},
	}

	board.Apply(0, 0, ToolFill)
	board.Apply(1, 0, ToolMark)
	board.Apply(0, 1, ToolMark)
	board.Apply(1, 1, ToolFill)

	if !IsSolved(board, solution) {
		t.Fatal("expected marked non-solution cells to be ignored for completion")
	}

	snapshot := board.Snapshot()
	board.ClearCell(1, 1)
	if IsSolved(board, solution) {
		t.Fatal("expected board to be unsolved after clearing a required filled cell")
	}
	board.Restore(snapshot)
	if !IsSolved(board, solution) {
		t.Fatal("expected restore to recover solved board")
	}

	if !board.SetCell(1, 1, CellEmpty) {
		t.Fatal("expected SetCell to clear a filled cell")
	}
	if board.SetCell(1, 1, CellEmpty) {
		t.Fatal("expected SetCell to report unchanged empty cell")
	}
}

func TestParseSolutionRejectsBadRows(t *testing.T) {
	p := Puzzle{
		Width:       3,
		Height:      2,
		SolutionRaw: []string{"010", "01"},
	}
	if err := p.ParseSolution(); err == nil {
		t.Fatal("expected width validation error")
	}
}

func assertClues(t *testing.T, got, want []int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("clue length got %v want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("clue got %v want %v", got, want)
		}
	}
}
