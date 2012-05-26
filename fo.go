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

func VerboseGetStat(player PlayerID, statname StatID, client StatsClient) {
	statval := client.GetStat(player, statname)
	log.Printf("%s, %v -> %f", player, statname, statval)
}

func FormatBattingStats(stats StatLine) string {
	return fmt.Sprintf("AVG:%0.3f HR:%02d R:%03d RBI:%03d SB:%02d",
		stats[B_BATTING_AVG],
		int(stats[B_HOME_RUNS]),
		int(stats[B_RUNS]),
		int(stats[B_RUNS_BATTED_IN]),
		int(stats[B_STOLEN_BASES]))
}

func FormatPitchingStats(stats StatLine) string {
	whip, ok := stats[P_WHIP]
	if !ok {
		whip = (stats[P_WALKS] + stats[P_HITS]) / stats[P_INNINGS]
	}

	return fmt.Sprintf("W:%02d S:%02d K:%03d ERA:%0.2f WHIP:%0.2f",
		int(stats[P_WINS]),
		0, // TODO: saves
		int(stats[P_STRIKE_OUTS]),
		stats[P_EARNED_RUN_AVERAGE],
		whip)
}

func (fo *FO) Optimize() {
	//	fo.zipsProjectMyRoster()
	//	fo.myCurrentStats()
	//	fo.leagueStats()

	rosters, err := fo.yahoo.LeagueRosters()
	if err != nil {
		log.Fatal(err)
	}

	for i := range *rosters {
		fmt.Printf("TEAM %d\n", i)
		starters := fo.selectStarters((*rosters)[i])
		for pos := range(starters) {
			fmt.Printf("%s ->  %v\n", pos, starters[pos])
		}
	}

//	for i := range *rosters {
//		fmt.Printf("TEAM %d\n", i)
//		fo.projectRoster((*rosters)[i], .9)
//	}

	// Full Docs:
	// http://developer.yahoo.com/fantasysports/guide/index.html
	//
	// League standings: 
	// "http://fantasysports.yahooapis.com/fantasy/v2/league/mlb.l.5181/standings",
	//
	// Team Roster:
	// "http://fantasysports.yahooapis.com/fantasy/v2/team/mlb.l.5181.t.6/roster",
	//
	// 10 Free Agents:
	// "http://fantasysports.yahooapis.com/fantasy/v2/league/mlb.l.5181/players;status=FA;count=10",
}

func (fo *FO) projectRoster(roster []YahooPlayer, seasonComplete float32) {
	for i := range roster {
		player := roster[i]
		stats := fo.projections.GetStatLine(PlayerID(player.FullName))

		if player.PositionType == "B" {
			fmt.Printf("%30s -> %s\n", player.FullName, FormatBattingStats(stats))
		} else {
			fmt.Printf("%30s -> %s\n", player.FullName, FormatPitchingStats(stats))
		}
	}
}

func (fo *FO) myCurrentStats() {
	mystats, err := fo.yahoo.MyStats()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n%s\n", FormatBattingStats(*mystats), FormatPitchingStats(*mystats))
}

func (fo *FO) leagueStats() {
	leaguestats, err := fo.yahoo.CurrentStats()
	if err != nil {
		log.Fatal(err)
	}

	for id, statline := range *leaguestats {
		fmt.Printf("%2d -> %s / %s\n", id, FormatBattingStats(statline), FormatPitchingStats(statline))
	}
}

func (fo *FO) zipsProjectMyRoster() {
	roster, err := fo.yahoo.MyRoster()
	if err != nil {
		log.Fatal(err)
	}

	for i := range *roster {
		name := (*roster)[i].FullName
		pType := (*roster)[i].PositionType
		stats := fo.projections.GetStatLine(PlayerID(name))
		if stats == nil {
			fmt.Printf("Couldn't get stats for '%s'\n", name)
		} else {
			if pType == "B" {
				fmt.Printf("%30s [%s] -> %s\n", name, pType, FormatBattingStats(stats))
			} else {
				fmt.Printf("%30s [%s] -> %s\n", name, pType, FormatPitchingStats(stats))
			}
		}
	}
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

// SimulateSeason
// - Fetch Current Stats
// - Fetch Rosters
// - Fetch Projections
// - Scale Projections
// - Compute Final Stats
