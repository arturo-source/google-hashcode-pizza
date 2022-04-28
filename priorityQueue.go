package main

type Node struct {
	p                Pizza
	currPos          int
	currPoints       int
	currSlices       Slices
	pessimisticBound int
}

type PriorityQueue struct {
	queue       []Node
	compareFunc func(introduced Node, inArray Node) bool
}

func (pq *PriorityQueue) Push(node Node) {
	low := 0
	high := len(pq.queue)
	for low < high {
		mid := (low + high) / 2

		if pq.compareFunc(node, pq.queue[mid]) {
			low = mid + 1
		} else {
			high = mid
		}
	}
	index := low

	if len(pq.queue) == index {
		pq.queue = append(pq.queue, node)
	} else {
		pq.queue = append(pq.queue[:index+1], pq.queue[index:]...)
		pq.queue[index] = node
	}
}

func (pq *PriorityQueue) Pop() Node {
	node := pq.queue[0]
	pq.queue = pq.queue[1:]

	return node
}

// func (pq *PriorityQueue) Print() {
// 	for _, n := range pq.queue {
// 		fmt.Print(n.pessimisticBound+n.currPoints, "-")
// 	}
// 	fmt.Println()
// }

func (pq *PriorityQueue) Empty() bool {
	return len(pq.queue) == 0
}

// Taken from https://rosettacode.org/wiki/Binary_search#Go
/* func binarySearch(a []int, value int) int {
	low := 0
	high := len(a) - 1
	for low <= high {
		mid := (low + high) / 2
		if a[mid] > value {
			high = mid - 1
		} else if a[mid] < value {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
} */

// Taken from https://stackoverflow.com/a/61822301
/* func insert(a []int, index int, value int) []int {
    if len(a) == index { // nil or empty slice or after last element
        return append(a, value)
    }
    a = append(a[:index+1], a[index:]...) // index < len(a)
    a[index] = value
    return a
} */
