package main

import (
	"fmt"
)

func greedyMethod(p Pizza, slices Slices, possiblePairs []PairType) (int, Slices) {
	var bestPoints int = 0
	var bestSlices Slices
	var lastSlice Slice
	if len(slices.slices) > 0 {
		lastSlice = slices.Last()
	}

	slice1x1 := Slice1x1(0, 0)

	for i := (lastSlice.y+lastSlice.height)*p.Columns + lastSlice.x + lastSlice.width; i < len(p.data); i++ {
		x, y := i%p.Columns, i/p.Columns
		slice1x1.x = x
		slice1x1.y = y

		if !slices.IsIn(slice1x1) {
			for _, pair := range possiblePairs {
				slice := Slice{
					x:      x,
					y:      y,
					width:  pair[0],
					height: pair[1],
				}

				if p.IsValidSlice(slice, slices) {
					pointsPerSlice := pair[0] * pair[1]

					bestPoints += pointsPerSlice
					bestSlices.slices = append(bestSlices.slices, slice)

					i += slice.width

					break
				}
			}
		}
	}

	return bestPoints, bestSlices
}

func optimisticBound(p Pizza, slices Slices) int {
	maxPoints := 0
	lastSlice := slices.Last()

	slice1x1 := Slice1x1(0, 0)

	for i := (lastSlice.y+lastSlice.height)*p.Columns + lastSlice.x + lastSlice.width; i < len(p.data); i++ {
		slice1x1.x = i % p.Columns
		slice1x1.y = i / p.Columns

		if !slices.IsIn(slice1x1) {
			maxPoints++
		}
	}

	return maxPoints
}

func BranchAndBound(p *Pizza, possiblePairs []PairType) (int, Slices) {
	// Pesimistic solution
	// var bestPoints int
	// var bestSlices Slices
	bestPoints, bestSlices := greedyMethod(*p, Slices{slices: []Slice{}}, possiblePairs)
	// Optimistic solution
	maxPoints := p.Rows * p.Columns

	pq := PriorityQueue{
		queue: []Node{},
		compareFunc: func(introduced Node, inArray Node) bool {
			// return introduced.currPoints+introduced.pessimisticBound < inArray.currPoints+inArray.pessimisticBound
			return true
		},
	}
	pq.Push(Node{
		currPoints: 0,
		currSlices: Slices{},
		// optimisticBound:  maxPoints,
		// pessimisticBound: bestPoints,
	})

	for !pq.Empty() {
		node := pq.Pop()
		// fmt.Println("len(pq.queue)", len(pq.queue))
		// fmt.Printf("%+v\n", node.currSlices)

		i, j := p.NextFreePositionFrom(node.currSlices)
		if i == -1 || j == -1 {
			fmt.Println("Arrived at leaf.")
			if node.currPoints > bestPoints {
				// Improves the solution!
				fmt.Println("Improve the solution! At leaf.")
				bestPoints = node.currPoints
				bestSlices = node.currSlices
			}
			continue
		}

		for _, pair := range possiblePairs {
			slice := Slice{
				x:      i,
				y:      j,
				width:  pair[0],
				height: pair[1],
			}

			isValid := p.IsValidSlice(slice, node.currSlices)
			if isValid {
				// Is feasible
				nextNodePoints := node.currPoints + slice.width*slice.height
				nextNodeSlices := Slices{slices: append(node.currSlices.slices, slice)}

				if nextNodePoints == maxPoints {
					// It's the best possible solution!
					fmt.Println("Finished with nextNodePoints == maxPoints")
					return nextNodePoints, nextNodeSlices
				}

				// pesimistic, slices := greedyMethod(*p, nextNodeSlices, possiblePairs)
				// optimistic := optimisticBound(*p, nextNodeSlices)
				// if nextNodePoints+pesimistic > bestPoints {
				// 	// Improve the solution!
				// 	fmt.Println("Improve the solution! With greedy.")
				// 	bestPoints = nextNodePoints + pesimistic
				// 	bestSlices.slices = append(nextNodeSlices.slices, slices.slices...)
				// 	fmt.Println("New best points:", bestPoints)
				// }
				if nextNodePoints > bestPoints {
					// Improve the solution!
					// fmt.Println("Improve the solution! With greedy.")
					bestPoints = nextNodePoints
					bestSlices = nextNodeSlices
					// fmt.Println("New best points:", bestPoints)
				}

				// if nextNodePoints+optimistic > bestPoints {
				// Is promising
				pq.Push(Node{
					currPoints: nextNodePoints,
					currSlices: nextNodeSlices,
					// optimisticBound:  optimistic,
					// pessimisticBound: pesimistic,
				})
				// }

				// if optimistic == pesimistic {
				// 	// It's the best possible solution!
				// 	fmt.Println("Finished with optimistic == pesimistic")
				// 	return nextNodePoints + pesimistic, Slices{append(nextNodeSlices.slices, slices.slices...)}
				// }
			}
		}
	}
	fmt.Println("Finished")

	return bestPoints, bestSlices
}
