package main

import (
  "testing"
)

func TestTwoTeamsOneStat(t *testing.T) {
	stats := map[TeamID]StatLine {
	1: StatLine {B_HOME_RUNS: 10},
	2: StatLine {B_HOME_RUNS: 5},
	}

	score := score(stats, map[StatID]struct{} { B_HOME_RUNS: struct{}{} })

	if score[1] != 2 { t.Errorf("Team 1 should have 2 points, has: %d", score[1]) }
	if score[2] != 1 { t.Errorf("Team 2 should have 1 points, has: %d", score[2]) }
}

func TestTwoTeamsTwoStats(t *testing.T) {
	stats := map[TeamID]StatLine {
	1: StatLine {B_HOME_RUNS: 10, P_STRIKE_OUTS: 100},
	2: StatLine {B_HOME_RUNS: 5, P_STRIKE_OUTS: 50},
	}

	score := score(stats, map[StatID]struct{} {
	B_HOME_RUNS: struct{}{},
	P_STRIKE_OUTS: struct{}{},
	})

	if score[1] != 4 { t.Errorf("Team 1 should have 4 points, has: %d", score[1]) }
	if score[2] != 2 { t.Errorf("Team 2 should have 2 points, has: %d", score[2]) }
}

func TwoTeamsReverseStats(t *testing.T) {
	stats := map[TeamID]StatLine {
	1: StatLine {B_HOME_RUNS: 10, P_EARNED_RUN_AVERAGE: 1.00},
	2: StatLine {B_HOME_RUNS: 5, P_EARNED_RUN_AVERAGE: 5.00},
	}

	score := score(stats, map[StatID]struct{} {
	B_HOME_RUNS: struct{}{},
	P_STRIKE_OUTS: struct{}{},
	})

	if score[1] != 4 { t.Errorf("Team 1 should have 4 points, has: %d", score[1]) }
	if score[2] != 2 { t.Errorf("Team 2 should have 2 points, has: %d", score[2]) }
}
