package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Missing input or output file\nExample of use: ./hashcode-pizza inputfile.in outputfile.out")
	}
	inputfile := os.Args[1]
	outputfile := os.Args[2]

	defer elapsed("Solution")()
	p, err := NewPizza(inputfile)
	if err != nil {
		panic(err)
	}

	pairs := GetAllPairs(p.AtLeast*2, p.Highest)
	fmt.Println(pairs)

	points, solution := BackTracking(p, pairs, 0)
	fmt.Println("points:", points)
	fmt.Println("Max posible points:", p.Columns*p.Rows)

	err = WriteSol(solution, outputfile)
	if err != nil {
		panic(err)
	}
}
