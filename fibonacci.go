package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	. "github.com/avoronkov/fibonacci-game/common"
	cur "github.com/rthornton128/goncurses"
)

// 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377 ...

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
