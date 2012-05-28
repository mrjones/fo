package main

import (
	"sort"
)

func score(stats map[TeamID]StatLine) map[TeamID]int {
	rawScores := make(map[StatID]map[TeamID]int)

	for statid := range(scoringCategories()) {
		rawScores[statid] = scoreStat(stats, statid)
	}

	return flatten(rawScores)
}

func scoreStat(stats map[TeamID]StatLine, statid StatID) map[TeamID]int {
	scoremap := make(map[TeamID]int)

	slice := statSlice(stats, statid)
	sort.Sort(slice)
	for teamid, statline := range(stats) {
		target := float64(statline[statid])
		score := slice.Search(target) + 1
		if lowerIsBetter(statid) {
			score = 10 - score + 1
		}
		scoremap[teamid] = score
	}

	return scoremap
}

func statSlice(stats map[TeamID]StatLine, statid StatID) sort.Float64Slice {
	slice := make(sort.Float64Slice, len(stats))
	i := 0
	for _, statline := range(stats) {
		slice[i] = float64(statline[statid])
		i++
	}
	return slice
}

func flatten(stats map[StatID]map[TeamID]int) map[TeamID]int {
	result := make(map[TeamID]int)
	for statid := range(stats) {
		for teamid := range(stats[statid]) {
			result[teamid] += stats[statid][teamid]
		}
	}

	return result
}

