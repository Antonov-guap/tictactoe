package main

import (
	"errors"
	"strconv"
)

type Game struct {
	field  gameField
	turn   cell
	over   bool
	winner cell
}

type gameField [][]cell

type cell byte

const (
	empty cell = iota
	cross
	zero
)

func createGame() (g Game) {
	g.field = createField(3)
	g.turn = cross
	return
}

func createField(size int) (f gameField) {
	f = make([][]cell, size)
	for i := range f {
		f[i] = make([]cell, size)
	}
	return
}

func (g *Game) MakeTurn(x, y int) (err error) {
	if g.IsOver() {
		return errors.New("wrong turn. Game is already over")
	}
	if err = g.field.put(g.turn, x, y); err != nil {
		return
	}
	if g.turn == cross {
		g.turn = zero
	} else {
		g.turn = cross
	}
	return
}

func (f *gameField) put(symbol cell, x, y int) (err error) {
	if x > len(*f) || y > len(*f) || x < 0 || y < 0 {
		return errors.New("wrong turn. Max value for cell is " + strconv.Itoa(len(*f)))
	}
	if (*f)[y][x] != empty {
		return errors.New("wrong turn. Cell is not empty")
	}
	(*f)[y][x] = symbol
	return
}

func (g *Game) IsOver() bool { // Empty if not finished
	if g.over {
		return true
	}

	for _, v := range []cell{cross, zero} {

		// diagonals
		toleft, toright := true, true

		for x := 0; x < len(g.field); x++ {

			// lanes
			cols, rows := true, true

			for y := 0; y < len(g.field) && (cols || rows); y++ {
				cols = cols && g.field[x][y] == v
				rows = rows && g.field[y][x] == v
			}

			toright = toright && g.field[x][x] == v
			toleft = toleft && g.field[len(g.field)-x-1][x] == v

			if cols || rows {
				g.over = true
				g.winner = v
				return true
			}
		}

		if toright || toleft {
			g.over = true
			g.winner = v
			return true
		}
	}

	var cnt int
	for i := range g.field {
		for j := range g.field[i] {
			if g.field[i][j] == empty {
				cnt++
			}
		}
	}

	return cnt == 0
}
