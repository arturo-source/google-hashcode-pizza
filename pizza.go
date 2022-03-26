package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Cell rune

const (
	Mushroom  Cell = 'M'
	Tomato    Cell = 'T'
	UsedCell  Cell = '*'
	ErrorCell Cell = '-'
)

type Pizza struct {
	Rows    int
	Columns int
	AtLeast int
	Highest int
	data    []Cell
}

func NewPizza(path string) (p Pizza, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	p.Rows, _ = strconv.Atoi(scanner.Text())
	scanner.Scan()
	p.Columns, _ = strconv.Atoi(scanner.Text())
	scanner.Scan()
	p.AtLeast, _ = strconv.Atoi(scanner.Text())
	scanner.Scan()
	p.Highest, _ = strconv.Atoi(scanner.Text())

	p.data = make([]Cell, p.Rows*p.Columns)

	for i := 0; i < p.Rows; i++ {
		scanner.Scan()
		row := scanner.Text()

		for j, cellType := range row {
			p.data[i*p.Columns+j] = Cell(cellType)
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return p, nil
}

func (p *Pizza) At(i, j int) (Cell, error) {
	if i >= p.Columns || j >= p.Rows {
		return ErrorCell, fmt.Errorf("Exceeded row or col. Got (i, j) (%d,%d) => Max is (%d, %d)", i, j, p.Columns, p.Rows)
	}
	return p.data[j*p.Columns+i], nil
}

func (p *Pizza) isUsedAt(i, j int) bool {
	// Assume that coordinates are always inside pizza
	cell, _ := p.At(i, j)
	return cell == UsedCell
}

func (p *Pizza) GetPizzaWithUsedSlice(slice *Slice) (pizzaCopy Pizza) {
	pizzaCopy = p.GetCopy()
	pizzaCopy.SetSliceUsed(slice)

	return
}

func (p *Pizza) GetCopy() (pizzaCopy Pizza) {
	pizzaCopy = *p
	pizzaCopy.data = make([]Cell, len(p.data))
	for i := range p.data {
		pizzaCopy.data[i] = p.data[i]
	}

	return
}

func (p *Pizza) SetSliceUsed(slice *Slice) {
	for i := slice.x; i < slice.x+slice.width; i++ {
		for j := slice.y; j < slice.y+slice.height; j++ {
			p.data[j*p.Columns+i] = UsedCell
		}
	}
}

func (p *Pizza) IsValidSlice(slice *Slice) bool {
	if slice.x+slice.width > p.Columns || slice.y+slice.height > p.Rows {
		return false
	}

	mushrooms, tomatoes := 0, 0

	for i := slice.x; i < slice.x+slice.width; i++ {
		for j := slice.y; j < slice.y+slice.height; j++ {
			if p.isUsedAt(i, j) {
				return false
			}

			cellContent, _ := p.At(i, j)
			switch cellContent {
			case Mushroom:
				mushrooms++
			case Tomato:
				tomatoes++
			default:
				fmt.Println("Something went wrong!!")
				return false
			}
		}
	}

	return mushrooms >= p.AtLeast && tomatoes >= p.AtLeast
}

func (p *Pizza) NextFreePositionFrom(currPos int) (i, j int) {
	for index := currPos; index < len(p.data); index++ {
		cell := p.data[index]
		if cell != UsedCell {
			j = index / p.Columns
			i = index - j*p.Columns
			return i, j
		}
	}
	// for j := 0; j < p.Rows; j++ {
	// 	for i := 0; i < p.Columns; i++ {
	// 		if !p.isUsedAt(i, j) {
	// 			return i, j
	// 		}
	// 	}
	// }

	return -1, -1
}

func (p *Pizza) Draw() {
	for j := 0; j < p.Rows; j++ {
		for i := 0; i < p.Columns; i++ {
			cell, _ := p.At(i, j)
			fmt.Print(string(cell))
		}
		fmt.Println()
	}
	fmt.Println()
}
