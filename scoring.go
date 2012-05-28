package main

import (
	"sort"
)

func score(stats map[TeamID]StatLine, scoringCategories map[StatID]struct{}) map[TeamID]float32 {
	scoresByStat := make(map[StatID]map[TeamID]float32)

	for statid := range scoringCategories {
		scoresByStat[statid] = scoreStat(stats, statid)
	}

	return flatten(scoresByStat)
}

func scoreStat(stats map[TeamID]StatLine, statid StatID) map[TeamID]float32 {
	numteams := len(stats)
	scoremap := make(map[TeamID]float32)

	slice := statSlice(stats, statid)
	sort.Sort(slice)
	for teamid, statline := range stats {
		target := float64(statline[statid])
		idx := slice.Search(target)
		score := float32(idx + 1)

		// TODO(mrjones): more generic tiebreaking
		if idx < numteams - 1 && slice[idx] == slice[idx + 1] {
			score += .5
		}
		if lowerIsBetter(statid) {
			score = float32(numteams) - score + 1
		}
		scoremap[teamid] = score
	}

	return scoremap
}

func statSlice(stats map[TeamID]StatLine, statid StatID) sort.Float64Slice {
	numteams := len(stats)
	slice := make(sort.Float64Slice, numteams)
	i := 0
	for _, statline := range stats {
		slice[i] = float64(statline[statid])
		i++
	}
	return slice
}

func flatten(stats map[StatID]map[TeamID]float32) map[TeamID]float32 {
	result := make(map[TeamID]float32)
	for statid := range stats {
		for teamid := range stats[statid] {
			result[teamid] += stats[statid][teamid]
		}
	}

	return result
}
