package folib

type ColIndex int
type ColName string
type PlayerID string
type Stat float64
type StatID int32

type Position string

type StatLine map[StatID]Stat

type TeamID int

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
)

func isRateStat(s StatID) bool {
	return s == B_BATTING_AVG ||
		s == B_ON_BASE_PCT ||
		s == B_SLUGGING ||
		s == P_EARNED_RUN_AVERAGE ||
		s == P_WHIP
}

func lowerIsBetter(s StatID) bool {
	return s == P_EARNED_RUN_AVERAGE ||
		s == P_WHIP
}

func merge(indiv []StatLine) StatLine {
	// replace equal-weight with unrolled/counting stats merge
	totals := make(StatLine)

	rawTotals := make(StatLine)
	counts := make(StatLine)

	for i := range indiv {
		for s := range indiv[i] {
			rawTotals[s] += indiv[i][s]
			if indiv[i][s] > 0.01 {
				counts[s] += 1
			}
		}
	}

	for s := range rawTotals {
		if isRateStat(s) {
			totals[s] = rawTotals[s] / counts[s]
		} else {
			totals[s] = rawTotals[s]
		}
	}

	return totals
}

type StatsClient interface {
	GetStat(player PlayerID, stat StatID) Stat
	GetStatLine(player PlayerID) StatLine
}
