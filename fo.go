package main

import (
	"fmt"
	"log"
	"sort"
)

const (
	LEAGUE_ID = 5181
)

type FO struct {
	yahoo       *YahooClient
	projections StatsClient
}

func NewFO(yahoo *YahooClient, projections StatsClient) *FO {
	return &FO{yahoo: yahoo, projections: projections}
}

func (fo *FO) Optimize() {
	rosters, err := fo.yahoo.LeagueRosters()
	if err != nil {
		log.Fatal(err)
	}

	teamStats, err := fo.yahoo.CurrentStats()
	if err != nil {
		log.Fatal(err)
	}

	teamProjections := make(map[TeamID]StatLine)
	for i := range *rosters {
		teamProjections[i] = fo.projectRoster((*rosters)[i], .9)
	}

	fmt.Printf("Projections\n")
	printScores(scoreLeague(teamProjections, scoringCategories()))

	fmt.Printf("\nActuals\n")
	printScores(scoreLeague(*teamStats, scoringCategories()))
}

func (fo *FO) projectPlayers(players []YahooPlayer, seasonComplete float32) map[PlayerID]StatLine {
	result := make(map[PlayerID]StatLine)

	for i := range(players) {
		id := PlayerID(players[i].FullName)
		result[id] = fo.projections.GetStatLine(id)
	}

	return result
}

func (fo *FO) projectRoster(roster []YahooPlayer, seasonComplete float32) StatLine {
	starterStats := make([]StatLine, 0)
	starters := fo.selectStarters(roster)

	for pos := range starters {
		for play := range starters[pos] {
			player := starters[pos][play]
			stats := fo.projections.GetStatLine(PlayerID(player.FullName))
			starterStats = append(starterStats, stats)
		}
	}

	return merge(starterStats)
}

type TeamLeaderEntry struct {
	Score float32
	ID PlayerID
}

type TeamLeaders []TeamLeaderEntry

func (l TeamLeaders) Len() int {
	return len(l)
}

func (l TeamLeaders) Less(i, j int) bool {
	return l[i].Score > l[j].Score
}

func (l TeamLeaders) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func SortedLeaders(scores map[PlayerID]float32) TeamLeaders {
	l := make(TeamLeaders, 0)
	for pid, score := range(scores) {
		l = append(l, TeamLeaderEntry{Score: score, ID: pid})
	}
	sort.Sort(l)
	return l
}

func indexByName(players []YahooPlayer) map[PlayerID]YahooPlayer {
	index := make(map[PlayerID]YahooPlayer)
	for _, player := range(players) {
		index[PlayerID(player.FullName)] = player
	}
	return index
}

func (fo *FO) selectStarters(roster []YahooPlayer) map[Position][]YahooPlayer {
	positionCounts := rosterTopology()
	statMap := fo.projectPlayers(roster, 1.0)
	leaders := SortedLeaders(scoreTeam(statMap, scoringCategories()))
	starters := make(map[Position][]YahooPlayer)
	index := indexByName(roster)

//	for i := range roster {
//		for j := range roster[i].Position {
//			pos := Position(roster[i].Position[j])
//			if positionCounts[pos] > 0 {
//				starters[pos] = append(starters[pos], roster[i])
//				positionCounts[pos]--
//			}
//			break
//		}
//	}

	for _, entry := range(leaders) {
		player := index[entry.ID]
		starting := false
		for _, posStr := range(player.Position) {
			pos := Position(posStr)
			if positionCounts[pos] > 0 {
				starters[pos] = append(starters[pos], player)
				positionCounts[pos]--
				fmt.Printf("%s is starting at %s\n", player.FullName, pos)
				starting = true
				break
			}
		}
		if !starting {
			fmt.Printf("%s is NOT starting\n", player.FullName)
		}
	}

	fmt.Println("---")

	return starters
}

type Position string

func rosterTopology() map[Position]int {
	return map[Position]int{
		"C":    1,
		"1B":   1,
		"2B":   1,
		"3B":   1,
		"SS":   1,
		"OF":   3,
		"Util": 3,
		"SP":   4,
		"RP":   2,
		"P":    2,
	}
}

func scoringCategories() map[StatID]struct{} {
	return map[StatID]struct{}{
		B_BATTING_AVG:    struct{}{},
		B_HOME_RUNS:      struct{}{},
		B_RUNS:           struct{}{},
		B_RUNS_BATTED_IN: struct{}{},
		B_STOLEN_BASES:   struct{}{},

		P_WINS:               struct{}{},
		P_SAVES:              struct{}{},
		P_EARNED_RUN_AVERAGE: struct{}{},
		P_WHIP:               struct{}{},
		P_STRIKE_OUTS:        struct{}{},
	}
}

// SimulateSeason
// - Fetch Current Stats
// - Fetch Rosters
// - Fetch Projections
// - Scale Projections
// - Compute Final Stats
