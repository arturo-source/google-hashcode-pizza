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

func WriteSol(slices []Slice, path string) (err error) {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(slices)))
	sb.WriteRune('\n')

	for i := range slices {
		slice := slices[len(slices)-1-i]
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", slice.y, slice.x, slice.y+slice.height-1, slice.x+slice.width-1))
	}

	err = os.WriteFile(path, []byte(sb.String()), 0644)
	return
}

func firstGreedyMethod(p Pizza, currPos int, pair PairType) (int, []Slice) {
	p = p.GetCopy()
	var bestPoints int = 0
	var bestSlices []Slice
	var pointsPerSlice int = pair[0] * pair[1]

	for i := currPos; i < len(p.data); i++ {
		j := i / p.Columns
		slice := Slice{
			x:      i % p.Columns,
			y:      j,
			width:  pair[0],
			height: pair[1],
		}
		if p.IsValidSlice(&slice) {
			p.SetSliceUsed(&slice)
			bestPoints += pointsPerSlice
			bestSlices = append(bestSlices, slice)
		}
	}

	return bestPoints, bestSlices
}

func greedyMethod(p Pizza, currPos int, possiblePairs []PairType) (int, []Slice) {
	p = p.GetCopy()
	var bestPoints int = 0
	var bestSlices []Slice

	for i := currPos; i < len(p.data); i++ {
		if p.data[i] != UsedCell {
			for _, pair := range possiblePairs {
				j := i / p.Columns
				slice := Slice{
					x:      i % p.Columns,
					y:      j,
					width:  pair[0],
					height: pair[1],
				}

				if p.IsValidSlice(&slice) {
					pointsPerSlice := pair[0] * pair[1]

					p.SetSliceUsed(&slice)
					bestPoints += pointsPerSlice
					bestSlices = append(bestSlices, slice)

					break
				}
			}
		}
	}

	return bestPoints, bestSlices
}

func optimisticBound(p Pizza, currPos int) int {
	maxPoints := 0

	for i := currPos; i < len(p.data); i++ {
		if p.data[i] != UsedCell {
			maxPoints++
		}
	}

	return maxPoints
}

func agressiveOptimisticBound(p Pizza, currPos int) int {
	// It's not valid
	tomatoPoints := 0
	mushroomPoints := 0

	for i := currPos; i < len(p.data); i++ {
		if p.data[i] == Tomato {
			tomatoPoints++
		} else if p.data[i] == Mushroom {
			mushroomPoints++
		}
	}

	if tomatoPoints > mushroomPoints {
		return mushroomPoints * 2
	} else {
		return tomatoPoints * 2
	}
}

func BackTracking(p Pizza, possiblePairs []PairType, currPos int) (int, []Slice) {
	i, j := p.NextFreePositionFrom(currPos)
	if i == -1 || j == -1 {
		// Is leaf
		return 0, nil
	}

	// Pesimistic solution
	bestPoints, bestSlices := greedyMethod(p, currPos, possiblePairs)
	for _, pair := range possiblePairs {
		slice := Slice{
			x:      i,
			y:      j,
			width:  pair[0],
			height: pair[1],
		}
		isValid := p.IsValidSlice(&slice)
		if !isValid {
			// Is not feasible
			continue
		}
		if optimisticBound(p, currPos) <= bestPoints {
			// Is not promising
			continue
		}

		pointsEarned, slices := BackTracking(p.GetPizzaWithUsedSlice(&slice), possiblePairs, currPos+slice.width)
		pairPoints := pair[0] * pair[1]
		if pairPoints+pointsEarned >= bestPoints {
			bestPoints = pairPoints + pointsEarned
			bestSlices = append(slices, slice)
		}
		if bestPoints >= p.Rows*p.Columns {
			return bestPoints, bestSlices
		}
	}

	return bestPoints, bestSlices
}

// Divide and conquer
func BruteForcing(p Pizza, possiblePairs []PairType, currPos int) (int, []Slice) {
	i, j := p.NextFreePositionFrom(currPos)
	if i == -1 || j == -1 {
		// Is leaf
		return 0, nil
	}

	var bestPoints int = 0
	var bestSlices []Slice
	for _, pair := range possiblePairs {
		slice := Slice{
			x:      i,
			y:      j,
			width:  pair[0],
			height: pair[1],
		}
		isValid := p.IsValidSlice(&slice)
		if !isValid {
			continue
		}

		pointsEarned, slices := BruteForcing(p.GetPizzaWithUsedSlice(&slice), possiblePairs, currPos+slice.width)
		pairPoints := pair[0] * pair[1]
		if pairPoints+pointsEarned >= bestPoints {
			bestPoints = pairPoints + pointsEarned
			bestSlices = append(slices, slice)
		}
	}

	return bestPoints, bestSlices
}

func BranchAndBound(p *Pizza, possiblePairs []PairType) (int, []Slice) {
	// Pesimistic solution
	bestPoints, bestSlices := greedyMethod(*p, 0, possiblePairs)
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
		currSlices:       []Slice{},
		optimisticBound:  maxPoints,
		pessimisticBound: bestPoints,
	})

	for !pq.Empty() {
		node := pq.Pop()
		p = node.p
		// fmt.Println("len(pq.queue)", len(pq.queue))

		i, j := p.NextFreePositionFrom(node.currPos)
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

			isValid := p.IsValidSlice(&slice)
			if isValid {
				// Is feasible
				nextNodePoints := node.currPoints + slice.width*slice.height
				nextNodeSlices := append(node.currSlices, slice)
				nextNodePos := node.currPos + slice.width

				if nextNodePoints == maxPoints {
					// It's the best possible solution!
					fmt.Println("Finished with nextNodePoints == maxPoints")
					return nextNodePoints, nextNodeSlices
				}

				pizzaUsed := p.GetPizzaWithUsedSlice(&slice)
				pesimistic, slices := greedyMethod(pizzaUsed, nextNodePos, possiblePairs)
				optimistic := optimisticBound(pizzaUsed, nextNodePos)
				if nextNodePoints+pesimistic > bestPoints {
					// Improve the solution!
					// fmt.Println("Improve the solution! With greedy.")
					bestPoints = nextNodePoints + pesimistic
					bestSlices = append(nextNodeSlices, slices...)
					// fmt.Println("New best points:", bestPoints)
				}

				if nextNodePoints+optimistic > bestPoints {
					// Is promising
					pq.Push(Node{
						p:                &pizzaUsed,
						currPos:          nextNodePos,
						currPoints:       nextNodePoints,
						currSlices:       nextNodeSlices,
						optimisticBound:  optimistic,
						pessimisticBound: pesimistic,
					})
				}

				if optimistic == pesimistic {
					// It's the best possible solution!
					fmt.Println("Finished with optimistic == pesimistic")
					return nextNodePoints + pesimistic, append(nextNodeSlices, slices...)
				}
			}
		}
	}
	fmt.Println("Finished")

	return bestPoints, bestSlices
}
