package common

import "math/rand"

type Direction int

const (
	Left  Direction = 1
	Right Direction = 2
	Up    Direction = 3
	Down  Direction = 4
)

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
	case Left:
		for y, line := range f.Data {
			for x, _ := range line {
				moved = f.movePoint(y, x, dir) || moved
			}
		}
	case Right:
		for y, _ := range f.Data {
			for x := 3; x >= 0; x-- {
				moved = f.movePoint(y, x, dir) || moved
			}
		}
	case Up:
		for x := 0; x < 4; x++ {
			for y := 0; y < 4; y++ {
				moved = f.movePoint(y, x, dir) || moved
			}
		}
	case Down:
		for x := 0; x < 4; x++ {
			for y := 3; y >= 0; y-- {
				moved = f.movePoint(y, x, dir) || moved
			}
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

func (f *Field) movePoint(y, x int, dir Direction) bool {
	cur := f.Data[y][x]
	if cur == 0 {
		return false
	}
	switch dir {
	case Left:
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
	case Right:
		for rx := x; rx < 3; rx++ {
			if pr := f.Data[y][rx+1]; pr != 0 {
				if next, ok := f.SumFib(pr, cur); ok {
					f.Data[y][rx+1], f.Data[y][x] = -next, 0
					f.score += next
					return true
				} else if rx != x {
					f.Data[y][rx], f.Data[y][x] = cur, 0
					return true
				}
				return false
			}
		}
		if x != 3 {
			f.Data[y][3], f.Data[y][x] = cur, 0
			return true
		}
	case Up:
		for uy := y; uy >= 1; uy-- {
			if pu := f.Data[uy-1][x]; pu != 0 {
				if next, ok := f.SumFib(pu, cur); ok {
					f.Data[uy-1][x], f.Data[y][x] = -next, 0
					f.score += next
					return true
				} else if uy != y {
					f.Data[uy][x], f.Data[y][x] = cur, 0
					return true
				}
				return false
			}
		}
		if y != 0 {
			f.Data[0][x], f.Data[y][x] = cur, 0
			return true
		}
	case Down:
		for dy := y; dy < 3; dy++ {
			if pd := f.Data[dy+1][x]; pd != 0 {
				if next, ok := f.SumFib(pd, cur); ok {
					f.Data[dy+1][x], f.Data[y][x] = -next, 0
					f.score += next
					return true
				} else if dy != y {
					f.Data[dy][x], f.Data[y][x] = cur, 0
					return true
				}
				return false
			}
		}
		if y != 3 {
			f.Data[3][x], f.Data[y][x] = cur, 0
			return true
		}
	default:
		panic(UnknownDirection{dir})
	}
	return false
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
