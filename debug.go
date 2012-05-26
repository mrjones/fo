package main

import (
	"fmt"
)

func print(scores map[TeamID]map[StatID]int, stats map[TeamID]StatLine) {
	for t := range(scores) {
		fmt.Printf("\nTEAM %d\n", t)
		for s := range(scores[t]) {
			fmt.Printf("Stat %d -> %f (%d)\n", s, stats[t][s], scores[t][s])
		}
	}
}
