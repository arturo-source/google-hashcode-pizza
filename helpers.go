package main

import (
	"fmt"
	"time"
)

type PairType [2]int

func GetAllPairs(moreThan, lessThan int) (pairs []PairType) {
	for n1 := 0; n1 <= lessThan; n1++ {
		for n2 := 0; n2 <= lessThan; n2++ {
			if n1*n2 >= moreThan {
				if n1*n2 <= lessThan {
					pairs = append(pairs, PairType{n1, n2})
				} else {
					break
				}
			}
		}
	}

	return
}

// Taken from https://stackoverflow.com/a/45766707
func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}
