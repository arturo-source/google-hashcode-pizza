package main

import (
	"fmt"
)

func greedyMethod(p Pizza, currPos int, possiblePairs []PairType) (int, Slices) {
	var bestPoints int = 0
	var bestSlices Slices

	pCopy := p.CopyPizza()

	for i := currPos; i < len(p.data); i++ {
		x, y := i%p.Columns, i/p.Columns

		if pCopy.data[i] != UsedCell {
			for _, pair := range possiblePairs {
				slice := Slice{
					x:      x,
					y:      y,
					width:  pair[0],
					height: pair[1],
				}

				if pCopy.IsValidSlice(slice) {
					pCopy.SetSliceUsed(slice)
					pointsPerSlice := pair[0] * pair[1]

					bestPoints += pointsPerSlice
					bestSlices.slices = append(bestSlices.slices, slice)

					i += slice.width - 1

					break
				}
			}
		}
	}

	return bestPoints, bestSlices
}

func BranchAndBound(p Pizza, possiblePairs []PairType) (int, Slices) {
	// Pesimistic solution
	bestPoints, bestSlices := greedyMethod(p, 0, possiblePairs)
	// Optimistic solution
	maxPoints := p.Rows * p.Columns

	pq := PriorityQueue{
		queue: []Node{},
		compareFunc: func(introduced Node, inArray Node) bool {
			return introduced.currPoints+introduced.pessimisticBound < inArray.currPoints+inArray.pessimisticBound
		},
	}
	pq.Push(Node{
		p:                p,
		currPos:          0,
		currPoints:       0,
		currSlices:       Slices{},
		pessimisticBound: bestPoints,
	})

	for !pq.Empty() {
		node := pq.Pop()
		p = node.p

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

			isValid := p.IsValidSlice(slice)
			if isValid {
				// Is feasible
				pizzaUsed := p.CopyPizza()
				pizzaUsed.SetSliceUsed(slice)

				nextNodePoints := node.currPoints + slice.width*slice.height
				nextNodeSlices := Slices{slices: append(node.currSlices.slices, slice)}
				nextNodePos := node.currPos + slice.width

				if nextNodePoints == maxPoints {
					// It's the best possible solution!
					fmt.Println("Finished with nextNodePoints == maxPoints")
					return nextNodePoints, nextNodeSlices
				}

				pesimistic, slices := greedyMethod(pizzaUsed, node.currPos, possiblePairs)
				if nextNodePoints+pesimistic > bestPoints {
					// Improve the solution!
					fmt.Println("Improve the solution! With greedy.")
					bestPoints = nextNodePoints + pesimistic
					bestSlices.slices = append(nextNodeSlices.slices, slices.slices...)
					fmt.Println("New best points:", bestPoints)
				}

				// Node is always promising
				// because I didn't find any valid optimistic heuristic
				pq.Push(Node{
					p:                pizzaUsed,
					currPos:          nextNodePos,
					currPoints:       nextNodePoints,
					currSlices:       nextNodeSlices,
					pessimisticBound: pesimistic,
				})

				if nextNodePoints+pesimistic == maxPoints {
					// It's the best possible solution!
					fmt.Println("Finished with nextNodePoints+pesimistic == maxPoints")
					return nextNodePoints + pesimistic, Slices{append(nextNodeSlices.slices, slices.slices...)}
				}
			}
		}
	}
	fmt.Println("Finished")

	return bestPoints, bestSlices
}
