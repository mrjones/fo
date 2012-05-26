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

	P_EARNED_RUNS        StatID = 1001
	P_EARNED_RUN_AVERAGE StatID = 1002
	P_GAMES              StatID = 1003
	P_HITS               StatID = 1004
	P_HOME_RUNS          StatID = 1005
	P_INNINGS            StatID = 1006
	P_LOSSES             StatID = 1007
	P_RUNS               StatID = 1008
	P_SAVES              StatID = 1009
	P_STARTS             StatID = 1010
	P_STRIKE_OUTS        StatID = 1011
	P_WALKS              StatID = 1012
	P_WHIP               StatID = 1013
	P_WINS               StatID = 1014

	map[StatID]bool RATE_STATS {
	B_BATTING_AVG: true,

	}
)


func (StatLine s1) add(StatLine s2, weight float) {
	for i := range(s2) {

	}
}

type StatsClient interface {
	GetStat(player PlayerID, stat StatID) Stat
	GetStatLine(player PlayerID) StatLine
}
