package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Slice struct {
	x, y, width, height int
}

func Slice1x1(x, y int) Slice {
	return Slice{
		x:      x,
		y:      y,
		width:  1,
		height: 1,
	}
}

func (s Slice) Draw() {
	for j := 0; j < s.height; j++ {
		for i := 0; i < s.width; i++ {
			fmt.Print("*")
		}
		fmt.Println()
	}
	fmt.Println()
}

type Slices struct {
	slices []Slice
}

func (s *Slices) WriteSol(path string) (err error) {
	slices := s.slices

	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(slices)))
	sb.WriteRune('\n')

	for _, slice := range slices {
		// slice := slices[len(slices)-1-i]
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", slice.y, slice.x, slice.y+slice.height-1, slice.x+slice.width-1))
	}

	err = os.WriteFile(path, []byte(sb.String()), 0644)
	return
}

func (s *Slices) IsIn(outSlice Slice) bool {
	for _, inSlice := range s.slices {
		if SlicesOverlap(inSlice, outSlice) {
			return true
		}
	}

	return false
}

func (s *Slices) Points() int {
	points := 0

	for _, slice := range s.slices {
		points += slice.width * slice.height
	}

	return points
}

func (s *Slices) Last() Slice {
	return s.slices[len(s.slices)-1]
}

// Taken from https://stackoverflow.com/a/306379
func SlicesOverlap(A, B Slice) bool {
	valueInRange := func(value, min, max int) bool {
		return (value > min) && (value < max)
	}

	xOverlap := valueInRange(A.x, B.x, B.x+B.width) ||
		valueInRange(B.x, A.x, A.x+A.width)

	yOverlap := valueInRange(A.y, B.y, B.y+B.height) ||
		valueInRange(B.y, A.y, A.y+A.height)

	return xOverlap && yOverlap
}
