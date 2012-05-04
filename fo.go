package main

import (
	"fmt"
	"log"
)

const (
	LEAGUE_ID = 5181
)

type FO struct {
	yahoo          *YahooClient
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
	whip := (stats[P_WALKS] + stats[P_HITS]) / stats[P_INNINGS]

	return fmt.Sprintf("W:%02d S:%02d K:%03d ERA:%0.2f WHIP:%0.2f",
		int(stats[P_WINS]),
		0, // TODO: saves
		int(stats[P_STRIKE_OUTS]),
		stats[P_EARNED_RUN_AVERAGE],
		whip)
}

func (fo *FO) Optimize() {
//	response, err := f.yahoo.Get(
//		"http://fantasysports.yahooapis.com/fantasy/v2/team/mlb.l.5181.t.6/roster")
//	if (err != nil) { log.Fatal(err) }
//	fmt.Println(response)

	fo.myCurrentStats()
//	fo.zipsProjectMyRoster()

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

func (fo *FO) myCurrentStats() {
//	resp, err := fo.yahoo.Get(
//		"http://fantasysports.yahooapis.com/fantasy/v2/league/mlb.l.5181/standings")
//	if (err != nil) { log.Fatal(err) }
//
//	fmt.Println(resp)

	fo.yahoo.MyStats()
}

func (fo *FO) zipsProjectMyRoster() {
	roster, err := fo.yahoo.MyRoster()
	if (err != nil) { log.Fatal(err) }

	for i := range *roster {
		name := (*roster)[i].FullName
		pType := (*roster)[i].PositionType
		stats := fo.projections.GetStatLine(PlayerID(name))
		if stats == nil {
			fmt.Printf("Couldn't get stats for '%s'\n", name)
		} else {
			if (pType == "B") {
				fmt.Printf("%30s [%s] -> %s\n", name, pType, FormatBattingStats(stats))
			} else {
				fmt.Printf("%30s [%s] -> %s\n", name, pType, FormatPitchingStats(stats))				
			}
		}
	}
}

// SimulateSeason
// - Fetch Current Stats
// - Fetch Rosters
// - Fetch Projections
// - Scale Projections
// - Compute Final Stats
