package common

import "math/rand"

type Direction int

const (
	// Direction should be positive (>0)
	_              = iota
	Left Direction = iota
	Right
	Up
	Down
)

func (d Direction) String() string {
	switch d {
	case Left:
		return "LEFT"
	case Right:
		return "RIGHT"
	case Up:
		return "UP"
	case Down:
		return "DOWN"
	default:
		return "UNKNOWN"
	}
}

type Field struct {
	Data     [][]int
	score    int
	sequence []int
}

func NewField() *Field {
	f := new(Field)
	f.Data = make([][]int, 4)
	for i, _ := range f.Data {
		f.Data[i] = make([]int, 4)
	}
	f.sequence = []int{1, 1, 2}
	return f
}

func (f *Field) Get(y, x int) int {
	return f.Data[y][x]
}

func (f *Field) Score() int {
	return f.score
}

func (f *Field) Sequence() []int {
	return f.sequence
}

func (f *Field) SumFib(a, b int) (next int, ok bool) {
	if a > b {
		a, b = b, a
	}
	for i := 0; i < len(f.sequence)-1; i++ {
		if f.sequence[i] == a && f.sequence[i+1] == b {
			if i == len(f.sequence)-2 {
				f.sequence = append(f.sequence, a+b)
			}
			next, ok = f.sequence[i+2], true
			break
		}
	}
	return
}

// Check if two numbers are next in Fibonacci sequence.
// Does not update sequence.
func (f *Field) FibNear(a, b int) bool {
	if a > b {
		a, b = b, a
	}
	for i := 0; i < len(f.sequence)-1; i++ {
		if f.sequence[i] == a && f.sequence[i+1] == b {
			return true
		}
	}
	return false
}

func (f *Field) AddPoint() bool {
	x := 1 + rand.Intn(2)
	empties := f.countEmptyCells()
	if empties == 0 {
		// game over
		return false
	}
	idx := rand.Intn(empties)
	f.fillEmptyCell(idx, x)
	return true
}

func (f *Field) Move(dir Direction) (moved bool) {
	switch dir {
	case Left, Up, Right, Down:
		if dir == Right {
			f.reflectVertically()
		}
		if dir == Down {
			f.reflectHorizontally()
		}
		if dir == Up || dir == Down {
			f.transpose()
		}
		for y, line := range f.Data {
			for x, _ := range line {
				moved = f.movePointLeft(y, x) || moved
			}
		}
		if dir == Up || dir == Down {
			f.transpose()
		}
		if dir == Down {
			f.reflectHorizontally()
		}
		if dir == Right {
			f.reflectVertically()
		}
	default:
		panic(UnknownDirection{dir})
	}

	for i, row := range f.Data {
		for j, _ := range row {
			if x := f.Data[i][j]; x < 0 {
				f.Data[i][j] = -x
			}
		}
	}

	return
}

func (f *Field) movePointLeft(y, x int) bool {
	cur := f.Data[y][x]
	if cur == 0 {
		return false
	}

	for lx := x; lx >= 1; lx-- {
		if pl := f.Data[y][lx-1]; pl != 0 {
			if next, ok := f.SumFib(pl, cur); ok {
				f.Data[y][lx-1], f.Data[y][x] = -next, 0
				f.score += next
				return true
			} else if lx != x {
				f.Data[y][lx], f.Data[y][x] = cur, 0
				return true
			}
			return false
		}
	}
	if x != 0 {
		f.Data[y][0], f.Data[y][x] = cur, 0
		return true
	}
	return false
}

func (f *Field) transpose() {
	for y, row := range f.Data {
		for x := y + 1; x < len(row); x++ {
			f.Data[x][y], f.Data[y][x] = f.Data[y][x], f.Data[x][y]
		}
	}
}

func (f *Field) reflectVertically() {
	for y, row := range f.Data {
		x1, x2 := 0, len(row)-1
		for x1 < x2 {
			f.Data[y][x1], f.Data[y][x2] = f.Data[y][x2], f.Data[y][x1]
			x1++
			x2--
		}
	}
}

func (f *Field) reflectHorizontally() {
	y1, y2 := 0, len(f.Data[0])-1
	for y1 < y2 {
		for x, _ := range f.Data[0] {
			f.Data[y1][x], f.Data[y2][x] = f.Data[y2][x], f.Data[y1][x]
		}
		y1++
		y2--
	}
}

// Check if there any possible moves. If no then game is over.
func (f *Field) HasPossibleMoves() bool {
	if f.countEmptyCells() > 0 {
		return true
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if j < 3 && f.FibNear(f.Data[i][j], f.Data[i][j+1]) {
				return true
			}
			if i < 3 && f.FibNear(f.Data[i][j], f.Data[i+1][j]) {
				return true
			}
		}
	}
	return false
}

func (f *Field) countEmptyCells() int {
	cnt := 0
	for _, line := range f.Data {
		for _, x := range line {
			if x == 0 {
				cnt++
			}
		}
	}
	return cnt
}

func (f *Field) fillEmptyCell(idx, value int) {
	for y, line := range f.Data {
		for x, content := range line {
			if content == 0 {
				if idx == 0 {
					f.Data[y][x] = value
					return
				}
				idx--
			}
		}
	}
}

type UnknownDirection struct {
	X Direction
}
