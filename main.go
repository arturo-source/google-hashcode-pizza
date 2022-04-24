package main

import (
	"fmt"
	"log"
	"os"
	"sort"
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
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][0]*pairs[i][1] > pairs[j][0]*pairs[j][1]
	})
	fmt.Println(pairs)

	points, solution := greedyMethod(p, 0, pairs)
	fmt.Println("greedy initial points:", points)

	points, solution = BranchAndBound(&p, pairs)
	fmt.Println("points:", points)
	fmt.Println("Max posible points:", p.Columns*p.Rows)

	err = WriteSol(solution, outputfile)
	if err != nil {
		panic(err)
	}
}
