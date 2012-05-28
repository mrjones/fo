package main

import (
	"fmt"
	"log"
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
	if err != nil { log.Fatal(err) }

	teamStats, err := fo.yahoo.CurrentStats()
	if err != nil { log.Fatal(err) }

	teamProjections := make(map[TeamID]StatLine)
	for i := range *rosters {
		teamProjections[i] = fo.projectRoster((*rosters)[i], .9)
	}

	fmt.Printf("Projections\n")
	printScores(score(teamProjections))

	fmt.Printf("\nActuals\n")
	printScores(score(*teamStats))
}

func (fo *FO) projectRoster(roster []YahooPlayer, seasonComplete float32) StatLine {
	starterStats := make([]StatLine, 0)
	starters := fo.selectStarters(roster)	

	for pos := range(starters) {
		for play := range(starters[pos]) {
			player := starters[pos][play]
			stats := fo.projections.GetStatLine(PlayerID(player.FullName))
			starterStats = append(starterStats, stats)
		}
	}

	return merge(starterStats)
}

func (fo *FO) selectStarters(roster []YahooPlayer) map[Position][]YahooPlayer {
	positionCounts := rosterTopology()
	starters := make(map[Position][]YahooPlayer)
	
	for i := range(roster) {
		for j := range(roster[i].Position) {
			pos := Position(roster[i].Position[j])
			if positionCounts[pos] > 0 {
				starters[pos] = append(starters[pos], roster[i])
				positionCounts[pos]--
			}
			break
		}
	}
	return starters
}

type Position string

func rosterTopology() map[Position]int {
	return map[Position]int {
		"C" : 1,
		"1B" : 1,
		"2B" : 1,
		"3B" : 1,
		"SS" : 1,
		"OF" : 3,
		"Util" : 3,
		"SP" : 4,
		"RP" : 2,
		"P" : 2,
	}
}

func scoringCategories() map[StatID]struct{} {
	return map[StatID]struct{} {
	B_BATTING_AVG: struct{}{},
	B_HOME_RUNS: struct{}{},
	B_RUNS: struct{}{},
	B_RUNS_BATTED_IN: struct{}{},
	B_STOLEN_BASES: struct{}{},

	P_WINS: struct{}{},
	P_SAVES: struct{}{},
	P_EARNED_RUN_AVERAGE: struct{}{},
	P_WHIP: struct{}{},
	P_STRIKE_OUTS: struct{}{},
	}
}

// SimulateSeason
// - Fetch Current Stats
// - Fetch Rosters
// - Fetch Projections
// - Scale Projections
// - Compute Final Stats
