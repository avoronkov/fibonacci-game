package common

import (
	"testing"
)

var testSequence = []int{1, 1, 2, 3, 5, 8}

func TestFieldMoveRight(t *testing.T) {
	before := Field{Data: [][]int{
		{0, 1, 0, 0},
		{0, 2, 3, 0},
		{1, 2, 3, 5},
		{0, 2, 2, 1},
	},
		sequence: testSequence,
	}
	before.Move(Right)
	after := [][]int{
		{0, 0, 0, 1},
		{0, 0, 0, 5},
		{0, 0, 3, 8},
		{0, 0, 2, 3},
	}
	for i, line := range before.Data {
		if exp := after[i]; !intSliceEqual(line, exp) {
			t.Errorf("Expected %v, found: %v", exp, line)
		}
	}
}

func TestMoveLeft(t *testing.T) {
	before := Field{Data: [][]int{
		{0, 1, 0, 0},
		{0, 2, 3, 0},
		{1, 2, 3, 5},
		{1, 2, 2, 0},
	},
		sequence: testSequence,
	}
	before.Move(Left)
	after := [][]int{
		{1, 0, 0, 0},
		{5, 0, 0, 0},
		{3, 8, 0, 0},
		{3, 2, 0, 0},
	}
	for i, line := range before.Data {
		if exp := after[i]; !intSliceEqual(line, exp) {
			t.Errorf("Expected %v, found %v", exp, line)
		}
	}
}

func TestMoveUp(t *testing.T) {
	before := Field{Data: [][]int{
		{0, 0, 1, 1},
		{0, 2, 2, 2},
		{1, 3, 3, 2},
		{0, 0, 5, 0},
	},
		sequence: testSequence,
	}
	before.Move(Up)
	after := [][]int{
		{1, 5, 3, 3},
		{0, 0, 8, 2},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	for i, line := range before.Data {
		if exp := after[i]; !intSliceEqual(line, exp) {
			t.Errorf("Expected %v, found %v", exp, line)
		}
	}
}

func TestMoveDown(t *testing.T) {
	before := Field{Data: [][]int{
		{0, 0, 1, 2},
		{0, 2, 2, 2},
		{1, 3, 3, 1},
		{0, 0, 5, 0},
	},
		sequence: testSequence,
	}
	before.Move(Down)
	after := [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 3, 2},
		{1, 5, 8, 3},
	}
	for i, line := range before.Data {
		if exp := after[i]; !intSliceEqual(line, exp) {
			t.Errorf("Expected %v, found %v", exp, line)
		}
	}
}

func intSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
