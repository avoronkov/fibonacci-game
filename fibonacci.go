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

var fib = []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377}

func SumFib(a, b int) (next int, ok bool) {
	if a > b {
		a, b = b, a
	}
	for i := 0; i < len(fib)-1; i++ {
		if fib[i] == a && fib[i+1] == b {
			next, ok = fib[i+2], true
			break
		}
	}
	return
}

type Field struct {
	Data  [][]int
	Score int
}

func NewField() *Field {
	f := new(Field)
	f.Data = make([][]int, 4)
	for i, _ := range f.Data {
		f.Data[i] = make([]int, 4)
	}
	return f
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
				if next, ok := SumFib(pl, cur); ok {
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
				if next, ok := SumFib(pr, cur); ok {
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
				if next, ok := SumFib(pu, cur); ok {
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
				if next, ok := SumFib(pd, cur); ok {
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
	colored   bool
	colorsMap []int16
}

func NewView(w *cur.Window) *View {
	v := new(View)
	v.win = w
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

func fibonacciIndex(x int) int16 {
	for i, val := range fib[1:] {
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
		v.win.AttrOn(cur.ColorPair(fibonacciIndex(x)))
		defer v.win.AttrOff(cur.ColorPair(fibonacciIndex(x)))
		if x >= 34 {
			v.win.AttrOn(cur.A_BOLD)
			defer v.win.AttrOff(cur.A_BOLD)
		}
	}
	v.win.Print(a)
}

func (v *View) drawLegend() {
	for _, f := range fib {
		v.colorPrint(f, f)
		v.win.Print("  ")
	}
	v.win.Print("...")
}

func (v *View) drawScore(score int) {
	v.win.Printf("Score: %v", score)
}

func (v *View) DrawField(f *Field) {
	height, width := v.win.MaxYX()
	startY, startX := height/2-5, width/2-15
	v.win.Move(startY-4, startX-13)
	v.drawLegend()
	for y, line := range f.Data {
		for x, value := range line {
			v.win.Move(startY+y, startX+5*x)
			v.colorPrint(value, formatNumber(value))
		}
	}
	v.win.Move(startY+10, startX)
	v.drawScore(f.Score)
}

func main() {
	win, err := cur.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer cur.End()

	cur.Echo(false)
	cur.Cursor(0)

	view := NewView(win)

	rand.Seed(time.Now().UTC().UnixNano())

	field := NewField()
	field.AddPoint()
	field.AddPoint()
L:
	for {
		view.DrawField(field)
		switch key := win.GetChar(); key {
		case cur.KEY_LEFT, 'h':
			if field.Move(Left) {
				field.AddPoint()
			}
		case cur.KEY_RIGHT, 'l':
			if field.Move(Right) {
				field.AddPoint()
			}
		case cur.KEY_UP, 'k':
			if field.Move(Up) {
				field.AddPoint()
			}
		case cur.KEY_DOWN, 'j':
			if field.Move(Down) {
				field.AddPoint()
			}
		case 'q':
			break L
		}
	}
}
