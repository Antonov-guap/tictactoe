package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var symbols = map[cell]rune{
	empty: ' ',
	cross: 'X',
	zero:  'O',
}

var xx = map[string]int{
	"A": 0,
	"B": 1,
	"C": 2,
}

func runStdout() {
	g := createGame()
	for !g.IsOver() {
		drawStdOut(g)
		makeTurnStdIn(&g)
	}
	drawStdOut(g)
}

func drawStdOut(g Game) {
	defer fmt.Println()

	for i, row := range g.field {
		for j, cell := range row {
			fmt.Printf(" %s", string(symbols[cell]))
			if j < len(row)-1 { // not last cell
				fmt.Print(" |")
			}
		}
		if i < len(g.field)-1 {
			fmt.Println()
			fmt.Println(strings.Repeat("-", len(row)*4-1))
		}
	}
	fmt.Println()

	if g.IsOver() {
		fmt.Printf("Game is over. \"%s\" wone.\n", string(symbols[g.winner]))
		return
	}
}

func makeTurnStdIn(g *Game) {
	var input string
	_, _ = fmt.Scanf("%s", &input)
	xy := strings.Split(input, "")
	if len(xy) != 2 {
		fmt.Printf("wrong input. Make a turn a1 or b2, for example")
		return
	}
	x := xx[strings.ToUpper(xy[0])]
	y, err := strconv.Atoi(xy[1])
	if err != nil {
		fmt.Printf("wrong input. Make a turn a1 or b2, for example")
		return
	}
	if err := g.MakeTurn(x, y-1); err != nil {
		log.Println(err)
	}
}
