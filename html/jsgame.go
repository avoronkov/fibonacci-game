package main

import (
	"math/rand"
	"time"

	game "github.com/avoronkov/fibonacci-game/common"
	"github.com/gopherjs/gopherjs/js"
)

func NewFieldJs() *js.Object {
	return js.MakeWrapper(game.NewField())
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	js.Global.Set("game", map[string]interface{}{
		"NewField": NewFieldJs,
		"Left":     game.Left,
		"Right":    game.Right,
		"Up":       game.Up,
		"Down":     game.Down,
	})
}
