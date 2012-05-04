package main

import (
	"fmt"
	"log"
)

const (
	LEAGUE_ID = 5181
)

type FO struct {
	yc *YahooClient
	fc *FangraphsClient
	projections *ZipsClient
}


func NewFO(yc *YahooClient, fc *FangraphsClient, projections *ZipsClient) *FO {
	return &FO{yc: yc, fc: fc, projections: projections}
}

func VerboseGetStat(player PlayerID, statname StatID, client *ZipsClient) {	
	statval := client.GetStat(player, statname)
	log.Printf("%s, %v -> %f", player, statname, statval)
}

func (f *FO) Optimize() {
//	response, err := f.yc.Get(
//		"http://fantasysports.yahooapis.com/fantasy/v2/team/mlb.l.5181.t.6/roster")
//	if (err != nil) { log.Fatal(err) }
//	fmt.Println(response)

//	f.fc.FetchAllData()
//	if err != nil { log.Fatal(err) }

//	fmt.Println(data)
//	if data != "" {
//		fmt.Println("done")
//	}

//	stat, err := f.projections.GetStat(PlayerID{FirstName: "Albert", LastName: "Pujols"}, BATTING_AVG)
	VerboseGetStat(PlayerID("Albert Pujols"), B_BATTING_AVG, f.projections)
VerboseGetStat(PlayerID("David Ortiz"), B_HOME_RUNS, f.projections)

kinsler := PlayerID("Ian Kinsler")
	VerboseGetStat(kinsler, B_BATTING_AVG, f.projections)
	VerboseGetStat(kinsler, B_RUNS, f.projections)
	VerboseGetStat(kinsler, B_RUNS_BATTED_IN, f.projections)
	VerboseGetStat(kinsler, B_HOME_RUNS, f.projections)
	VerboseGetStat(kinsler, B_STOLEN_BASES, f.projections)
	

	log.Print("Done");
	fmt.Println("Done");

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


// SimulateSeason
// - Fetch Current Stats
// - Fetch Rosters
// - Fetch Projections
// - Scale Projections
// - Compute Final Stats
