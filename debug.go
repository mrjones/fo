package main

import (
	"fmt"
)

func print(scores map[TeamID]map[StatID]int, stats map[TeamID]StatLine) {
	for t := range(scores) {
		fmt.Printf("\nTEAM %d\n", t)
		for s := range(scores[t]) {
			fmt.Printf("Stat %d -> %f (%d)\n", s, stats[t][s], scores[t][s])
		}
	}
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
