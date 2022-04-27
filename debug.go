package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func DrawCutPizza(pizzafile, slicesfile string) error {
	p, err := NewPizza(pizzafile)
	if err != nil {
		return err
	}
	slices, err := NewSlices(slicesfile)
	if err != nil {
		return err
	}

	for _, slice := range slices.slices {
		for i := slice.x; i < slice.x+slice.width; i++ {
			for j := slice.y; j < slice.y+slice.height; j++ {
				if p.data[j*p.Columns+i] == UsedCell {
					p.data[j*p.Columns+i] = ErrorCell
				} else {
					p.data[j*p.Columns+i] = UsedCell
				}
			}
		}
	}

	for j := 0; j < p.Rows; j++ {
		for i := 0; i < p.Columns; i++ {
			cell, _ := p.At(i, j)
			if cell != ErrorCell {
				fmt.Print(string(cell))
			} else {
				fmt.Print("\033[31m*\033[m")
			}
		}
		fmt.Println()
	}

	return nil
}

func NewSlices(path string) (slices Slices, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	rows, _ := strconv.Atoi(scanner.Text())

	slices.slices = make([]Slice, rows)

	for i := 0; i < rows; i++ {
		scanner.Scan()
		y := scanner.Text()
		scanner.Scan()
		x := scanner.Text()
		scanner.Scan()
		heigt := scanner.Text()
		scanner.Scan()
		width := scanner.Text()

		s := &slices.slices[i]

		s.x, _ = strconv.Atoi(x)
		s.y, _ = strconv.Atoi(y)
		s.width, _ = strconv.Atoi(width)
		s.height, _ = strconv.Atoi(heigt)
		s.width += 1 - s.x
		s.height += 1 - s.y
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}
