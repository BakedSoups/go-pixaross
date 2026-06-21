package nonogram

type CellState uint8

const (
	CellEmpty CellState = iota
	CellFilled
	CellMarked
)

type Tool uint8

const (
	ToolFill Tool = iota
	ToolMark
)

type Board struct {
	Width  int
	Height int
	Cells  [][]CellState
}

func NewBoard(width, height int) *Board {
	b := &Board{Width: width, Height: height}
	b.Cells = make([][]CellState, height)
	for y := range b.Cells {
		b.Cells[y] = make([]CellState, width)
	}
	return b
}

func (b *Board) Apply(x, y int, tool Tool) bool {
	return b.SetCell(x, y, TargetState(tool))
}

func (b *Board) SetCell(x, y int, state CellState) bool {
	if !b.InBounds(x, y) {
		return false
	}
	if b.Cells[y][x] == state {
		return false
	}
	b.Cells[y][x] = state
	return true
}

func (b *Board) ClearCell(x, y int) bool {
	if !b.InBounds(x, y) || b.Cells[y][x] == CellEmpty {
		return false
	}
	b.Cells[y][x] = CellEmpty
	return true
}

func (b *Board) Reset() {
	for y := range b.Cells {
		for x := range b.Cells[y] {
			b.Cells[y][x] = CellEmpty
		}
	}
}

func (b *Board) Snapshot() [][]CellState {
	copyCells := make([][]CellState, b.Height)
	for y := range b.Cells {
		copyCells[y] = append([]CellState(nil), b.Cells[y]...)
	}
	return copyCells
}

func (b *Board) Restore(snapshot [][]CellState) {
	for y := range b.Cells {
		copy(b.Cells[y], snapshot[y])
	}
}

func (b *Board) InBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < b.Width && y < b.Height
}

func TargetState(tool Tool) CellState {
	if tool == ToolMark {
		return CellMarked
	}
	return CellFilled
}
