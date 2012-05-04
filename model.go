package main

type ColIndex int
type ColName string
type PlayerID string
type Stat float64
type StatID int32

type StatLine map[StatID]Stat

const (
	B_AT_BATS         StatID = 1
	B_BATTING_AVG     StatID = 2
	B_CAUGHT_STEALING StatID = 3
	B_DOUBLES         StatID = 4
	B_GAMES           StatID = 5
	B_HITS            StatID = 6
	B_HOME_RUNS       StatID = 7
	B_ON_BASE_PCT     StatID = 8
	B_PLATE_APPS      StatID = 9
	B_RUNS            StatID = 10
	B_RUNS_BATTED_IN  StatID = 11
	B_SLUGGING        StatID = 12
	B_STOLEN_BASES    StatID = 13
	B_STRIKE_OUTS     StatID = 14
	B_TRIPLES         StatID = 15
	B_WALKS           StatID = 16
)
