package main

import (
	"fmt"
	"log"
)

type FO struct {
	yc *YahooClient
}

func NewFO(yc *YahooClient) *FO {
	return &FO{yc: yc}
}

func (f *FO) Optimize() {
	response, err := f.yc.Get(
		"http://fantasysports.yahooapis.com/fantasy/v2/team/mlb.l.5181.t.6/roster")
	if (err != nil) { log.Fatal(err) }
	fmt.Println(response)

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
