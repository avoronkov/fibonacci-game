package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	cur "github.com/rthornton128/goncurses"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

// 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377 ...

type Field struct {
	Data     [][]int
	Score    int
	Sequence []int
}

func NewField() *Field {
	f := new(Field)
	f.Data = make([][]int, 4)
	for i, _ := range f.Data {
		f.Data[i] = make([]int, 4)
	}
	f.Sequence = []int{1, 1, 2}
	return f
}

func (f *Field) SumFib(a, b int) (next int, ok bool) {
	if a > b {
		a, b = b, a
	}
	for i := 0; i < len(f.Sequence)-1; i++ {
		if f.Sequence[i] == a && f.Sequence[i+1] == b {
			if i == len(f.Sequence)-2 {
				f.Sequence = append(f.Sequence, a+b)
			}
			next, ok = f.Sequence[i+2], true
			break
		}
	}
	return
}

// Check if two numbers are next in Fibonacci sequence.
// Does not update Sequence.
func (f *Field) FibNear(a, b int) bool {
	if a > b {
		a, b = b, a
	}
	for i := 0; i < len(f.Sequence)-1; i++ {
		if f.Sequence[i] == a && f.Sequence[i+1] == b {
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
		panic(fmt.Errorf("Move(): unknown direction: %v", dir))
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
					f.Score += next
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
					f.Score += next
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
					f.Score += next
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
					f.Score += next
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
		panic(fmt.Errorf("Move(): unknown direction: %v", dir))
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

type View struct {
	win       *cur.Window
	field     *Field
	colored   bool
	colorsMap []int16
}

func NewView(w *cur.Window, f *Field) *View {
	v := new(View)
	v.win = w
	v.field = f
	v.colored = cur.HasColors()
	if v.colored {
		if err := cur.StartColor(); err != nil {
			log.Print(err)
		}
		if err := cur.UseDefaultColors(); err != nil {
			log.Print(err)
		}
		v.colorsMap = []int16{
			cur.C_WHITE,
			cur.C_YELLOW,
			cur.C_GREEN,
			cur.C_BLUE,
			cur.C_CYAN,
			cur.C_MAGENTA,
			cur.C_RED,
			cur.C_YELLOW,
			cur.C_GREEN,
			cur.C_BLUE,
			cur.C_CYAN,
			cur.C_MAGENTA,
			cur.C_RED,
		}
		for i, col := range v.colorsMap {
			if err := cur.InitPair(int16(i), col, cur.C_BLACK); err != nil {
				log.Printf("InitPair(%d, %v, %v) failed: %v", i, col, cur.C_BLACK, err)
			}
		}
	}
	return v
}

func (v *View) fibonacciIndex(x int) int16 {
	if x == 0 {
		return 0
	}
	for i, val := range v.field.Sequence[1:] {
		if x == val {
			return int16(i)
		}
	}
	log.Printf("Unknown fibonacci index: %v", x)
	return 0
}

func formatNumber(x int) string {
	if x == 0 {
		return "  .  "
	}
	s := fmt.Sprintf("%d", x)
	for len(s) < 5 {
		s = " " + s
		if len(s) >= 5 {
			return s
		}
		s = s + " "
	}
	return s
}

func (v *View) colorPrint(x int, a interface{}) {
	if v.colored {
		fi := v.fibonacciIndex(x)
		v.win.AttrOn(cur.ColorPair(fi))
		defer v.win.AttrOff(cur.ColorPair(fi))
		if x >= 34 {
			v.win.AttrOn(cur.A_BOLD)
			defer v.win.AttrOff(cur.A_BOLD)
		}
	}
	v.win.Print(a)
}

func (v *View) drawLegend(y int) {
	s := ""
	for _, f := range v.field.Sequence {
		s = s + fmt.Sprintf("%d  ", f)
	}
	s += "..."
	_, width := v.win.MaxYX()
	x := width/2 - len(s)/2
	v.win.Move(y, x)
	for _, f := range v.field.Sequence {
		v.colorPrint(f, f)
		v.win.Print("  ")
	}
	v.win.Print("...")
}

func (v *View) drawScore(score int) {
	v.win.Printf("Score: %v", score)
}

func (v *View) DrawField() {
	height, width := v.win.MaxYX()
	startY, startX := height/2-5, width/2-10
	// v.win.Move(startY-4, startX-13)
	v.drawLegend(startY - 4)
	for y, line := range v.field.Data {
		for x, value := range line {
			v.win.Move(startY+y, startX+5*x)
			v.colorPrint(value, formatNumber(value))
		}
	}
	v.win.Move(startY+7, startX+2)
	v.drawScore(v.field.Score)
}

func (v *View) GameOver() {
	height, width := v.win.MaxYX()
	startY, startX := height/2, width/2-8
	v.win.Move(startY, startX)
	v.win.Print("Game over.")
}

func main() {
	win, err := cur.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer cur.End()

	rand.Seed(time.Now().UTC().UnixNano())

	cur.Echo(false)
	cur.Cursor(0)

	field := NewField()
	field.AddPoint()
	field.AddPoint()

	view := NewView(win, field)

L:
	for {
		view.DrawField()
		if !field.HasPossibleMoves() {
			view.GameOver()
		}
		switch key := win.GetChar(); key {
		case cur.KEY_LEFT, 'h', 'a':
			if field.Move(Left) {
				field.AddPoint()
			}
		case cur.KEY_RIGHT, 'l', 'd':
			if field.Move(Right) {
				field.AddPoint()
			}
		case cur.KEY_UP, 'k', 'w':
			if field.Move(Up) {
				field.AddPoint()
			}
		case cur.KEY_DOWN, 'j', 's':
			if field.Move(Down) {
				field.AddPoint()
			}
		case 'q':
			break L
		}
	}
}
