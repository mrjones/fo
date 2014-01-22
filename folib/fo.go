package folib

import (
	"fmt"
	"log"
	"sort"
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
	if err != nil {
		log.Fatal(err)
	}

	err = fo.scoreTrade(rosters, "Matt Cain", "Troy Tulowitzki")
	if err != nil {
		log.Fatal(err)
	}

//	teamStats, err := fo.yahoo.CurrentStats()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	teamProjections := projectLeague(rosters)
//
//	fmt.Printf("Projections\n")
//	printScores(scoreLeague(teamProjections, scoringCategories()))
//
//	fmt.Printf("\nActuals\n")
//	printScores(scoreLeague(*teamStats, scoringCategories()))
}

func (fo *FO) projectLeague(rosters *map[TeamID][]YahooPlayer) map[TeamID]StatLine {
	teamProjections := make(map[TeamID]StatLine)
	for i := range *rosters {
		teamProjections[i] = fo.projectRoster((*rosters)[i], .9)
	}
	return teamProjections
}

func (fo *FO) scoreTrade(rosters *map[TeamID][]YahooPlayer, p1, p2 PlayerID) error {
	beforeProjections := fo.projectLeague(rosters)
	// copy?
		err, t1, t2 := trade(rosters, p1, p2)
	if err != nil {
		return err
	}
		err, t1, t2 = trade(rosters, "Zack Cozart", "Justin Masterson")
	if err != nil {
		return err
	}
	afterProjections := fo.projectLeague(rosters)

	fmt.Printf("Before\n")
	beforeScores := scoreLeague(beforeProjections, scoringCategories())
	printScores(beforeScores)
	fmt.Printf("TEAM %d: %s -> %s\n", t1, FormatBattingStats(beforeProjections[t1]), FormatBattingStats(afterProjections[t1]))
	fmt.Printf("TEAM %d: %s -> %s\n", t1, FormatPitchingStats(beforeProjections[t1]), FormatPitchingStats(afterProjections[t1]))
	fmt.Printf("TEAM %d: %s -> %s\n", t2, FormatBattingStats(beforeProjections[t2]), FormatBattingStats(afterProjections[t2]))
	fmt.Printf("TEAM %d: %s -> %s\n", t2, FormatPitchingStats(beforeProjections[t2]), FormatPitchingStats(afterProjections[t2]))

	fmt.Printf("After\n")
	afterScores := scoreLeague(afterProjections, scoringCategories())
	printScores(afterScores)

	fmt.Printf("Delta\n")
	for t := range(beforeProjections) {
		fmt.Printf("TEAM %d: %f\n", t, afterScores[t] - beforeScores[t])
	}

	return nil
}

func trade(rosters *map[TeamID][]YahooPlayer, p1, p2 PlayerID) (error, TeamID, TeamID) {
	t1, i1 := TeamID(-1), -1
	t2, i2 := TeamID(-1), -1
	for team, roster := range(*rosters) {
		for i, player := range(roster) {
			c := PlayerID(player.FullName)
			if p1 == c {
				t1, i1 = team, i
			} else if (p2 == c) {
				t2, i2 = team, i
			}
		}
	}

	if t1 == -1 || i1 == -1 || t2 == -1 || i2 == 1 {
		return fmt.Errorf("Bad lookup: %d %d %d %d", t1, i1, t2, i2), -1, -1
	}

	fmt.Printf("Moving %s from team %d to team %d\n", p1, t1, t2)
	fmt.Printf("Moving %s from team %d to team %d\n", p2, t2, t1)

	(*rosters)[t1][i1], (*rosters)[t2][i2] = (*rosters)[t2][i2], (*rosters)[t1][i1]
	return nil, t1, t2
}

func (fo *FO) projectPlayers(players []YahooPlayer, seasonComplete float32) map[PlayerID]StatLine {
	result := make(map[PlayerID]StatLine)

	for i := range players {
		id := PlayerID(players[i].FullName)
		result[id] = fo.projections.GetStatLine(id)
	}

	return result
}

func (fo *FO) projectRoster(roster []YahooPlayer, seasonComplete float32) StatLine {
	starterStats := make([]StatLine, 0)
	starters := fo.selectStarters(roster)

	for pos := range starters {
		for play := range starters[pos] {
			player := starters[pos][play]
			stats := fo.projections.GetStatLine(PlayerID(player.FullName))
			starterStats = append(starterStats, stats)
		}
	}

	return merge(starterStats)
}

type TeamLeaderEntry struct {
	Score float32
	ID    PlayerID
}

type TeamLeaders []TeamLeaderEntry

func (l TeamLeaders) Len() int {
	return len(l)
}

func (l TeamLeaders) Less(i, j int) bool {
	return l[i].Score > l[j].Score
}

func (l TeamLeaders) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func SortedLeaders(scores map[PlayerID]float32) TeamLeaders {
	l := make(TeamLeaders, 0)
	for pid, score := range scores {
		l = append(l, TeamLeaderEntry{Score: score, ID: pid})
	}
	sort.Sort(l)
	return l
}

func indexByName(players []YahooPlayer) map[PlayerID]YahooPlayer {
	index := make(map[PlayerID]YahooPlayer)
	for _, player := range players {
		index[PlayerID(player.FullName)] = player
	}
	return index
}

func (fo *FO) selectStarters(roster []YahooPlayer) map[Position][]YahooPlayer {
	positionCounts := rosterTopology()
	statMap := fo.projectPlayers(roster, 1.0)
	leaders := SortedLeaders(scoreTeam(statMap, scoringCategories()))
	starters := make(map[Position][]YahooPlayer)
	index := indexByName(roster)

	for _, entry := range leaders {
		player := index[entry.ID]
		starting := false
		for _, posStr := range player.Position {
			pos := Position(posStr)
			if positionCounts[pos] > 0 {
				starters[pos] = append(starters[pos], player)
				positionCounts[pos]--
//				fmt.Printf("%s is starting at %s\n", player.FullName, pos)
				starting = true
				break
			}
		}
		if !starting {
//			fmt.Printf("%s is NOT starting\n", player.FullName)
		}
	}

//	fmt.Println("---")

	return starters
}

func rosterTopology() map[Position]int {
	return map[Position]int{
		"C":    1,
		"1B":   1,
		"2B":   1,
		"3B":   1,
		"SS":   1,
		"OF":   3,
		"Util": 3,
		"SP":   4,
		"RP":   2,
		"P":    2,
	}
}

func scoringCategories() map[StatID]struct{} {
	return map[StatID]struct{}{
		B_BATTING_AVG:    struct{}{},
		B_HOME_RUNS:      struct{}{},
		B_RUNS:           struct{}{},
		B_RUNS_BATTED_IN: struct{}{},
		B_STOLEN_BASES:   struct{}{},

		P_WINS:               struct{}{},
		P_SAVES:              struct{}{},
		P_EARNED_RUN_AVERAGE: struct{}{},
		P_WHIP:               struct{}{},
		P_STRIKE_OUTS:        struct{}{},
	}
}

// SimulateSeason
// - Fetch Current Stats
// - Fetch Rosters
// - Fetch Projections
// - Scale Projections
// - Compute Final Stats
