package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

type ZipsClient struct {
	battingStats *map[PlayerID]StatLine
}

func (zc *ZipsClient) GetStat(player PlayerID, stat StatID) Stat {
	return (*zc.battingStats)[player][stat]
}

func (zc *ZipsClient) GetStatLine(player PlayerID) StatLine {
	return (*zc.battingStats)[player]
}

func NewZipsClient() (*ZipsClient, error) {
	EnsureCache(
		"http://www.baseballthinkfactory.org/szymborski/ZiPS2012v1BAT.csv",
		"zips2012batters.csv")
	EnsureCache(
		"http://www.baseballthinkfactory.org/szymborski/ZiPS2012v1PIT.csv",
		"zips2012pitchers.csv")

	battingStats, err := indexBattingStats()
	if err != nil {
		return nil, err
	}

	return &ZipsClient{battingStats: battingStats}, nil
}

func mapColumnNameToStat() map[ColName]StatID {
	return map[ColName]StatID{
		"AB":  B_AT_BATS,
		"BA":  B_BATTING_AVG,
		"CS":  B_CAUGHT_STEALING,
		"2B":  B_DOUBLES,
		"G":   B_GAMES,
		"H":   B_HITS,
		"HR":  B_HOME_RUNS,
		"OBP": B_ON_BASE_PCT,
		"PA":  B_PLATE_APPS,
		"R":   B_RUNS,
		"RBI": B_RUNS_BATTED_IN,
		"SLG": B_SLUGGING,
		"SB":  B_STOLEN_BASES,
		"K":   B_STRIKE_OUTS,
		"3B":  B_TRIPLES,
		"BB":  B_WALKS,
	}
}

func mapStatToColumnIndex(colNames []string) map[StatID]ColIndex {
	statToColumnIndex := map[StatID]ColIndex{}
	columnNameToStat := mapColumnNameToStat()
	for i := range colNames {
		colName := ColName(colNames[i])
		if statId, ok := columnNameToStat[colName]; ok {
			statToColumnIndex[statId] = ColIndex(i)
		}
	}

	return statToColumnIndex
}

func indexBattingStats() (*map[PlayerID]StatLine, error) {
	f, err := os.Open("zips2012batters.csv")
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	recs, err := r.ReadAll()

	if err != nil {
		return nil, err
	}

	statToColumnIndex := mapStatToColumnIndex(recs[0])

	statIndex := map[PlayerID]StatLine{}
	for i := 1; i < len(recs); i++ {
		statLine := make(StatLine)
		for stat, index := range statToColumnIndex {
			statStr := recs[i][index]
			stat64, err := strconv.ParseFloat(statStr, 64)
			if err != nil {
				return nil, err
			}
			statLine[stat] = Stat(stat64)
		}
		statIndex[PlayerID(recs[i][0])] = statLine
	}

	return &statIndex, nil
}
