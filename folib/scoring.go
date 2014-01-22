package folib

import (
	"sort"
	"strconv"
)

func scoreLeague(stats map[TeamID]StatLine, scoringCategories map[StatID]struct{}) map[TeamID]float32 {
	rawStats := make(map[string]StatLine)
	for k, v := range stats {
		rawStats[strconv.Itoa(int(k))] = v
	}

	rawScores := score(rawStats, scoringCategories)
	scores := make(map[TeamID]float32)
	for k, v := range rawScores {
		i, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		scores[TeamID(i)] = v
	}

	return scores
}

func scoreTeam(stats map[PlayerID]StatLine, scoringCategories map[StatID]struct{}) map[PlayerID]float32 {
	rawStats := make(map[string]StatLine)
	for k, v := range stats {
		rawStats[string(k)] = v
	}

	rawScores := score(rawStats, scoringCategories)
	scores := make(map[PlayerID]float32)
	for k, v := range rawScores {
		scores[PlayerID(k)] = v
	}

	return scores
}

func score(stats map[string]StatLine, scoringCategories map[StatID]struct{}) map[string]float32 {
	scoresByStat := make(map[StatID]map[string]float32)

	for statid := range scoringCategories {
		scoresByStat[statid] = scoreStat(stats, statid)
	}

	return flatten(scoresByStat)
}

func scoreStat(stats map[string]StatLine, statid StatID) map[string]float32 {
	numteams := len(stats)
	scoremap := make(map[string]float32)

	slice := statSlice(stats, statid)
	sort.Sort(slice)
	for teamid, statline := range stats {
		target := float64(statline[statid])
		idx := slice.Search(target)
		score := float32(idx + 1)

		// TODO(mrjones): more generic tiebreaking
		if idx < numteams-1 && slice[idx] == slice[idx+1] {
			score += .5
		}
		if lowerIsBetter(statid) {
			score = float32(numteams) - score + 1
		}
		scoremap[teamid] = score
	}

	return scoremap
}

func statSlice(stats map[string]StatLine, statid StatID) sort.Float64Slice {
	numteams := len(stats)
	slice := make(sort.Float64Slice, numteams)
	i := 0
	for _, statline := range stats {
		slice[i] = float64(statline[statid])
		i++
	}
	return slice
}

func flatten(stats map[StatID]map[string]float32) map[string]float32 {
	result := make(map[string]float32)
	for statid := range stats {
		for teamid := range stats[statid] {
			result[teamid] += stats[statid][teamid]
		}
	}

	return result
}
