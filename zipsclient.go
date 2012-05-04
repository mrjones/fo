package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

const (
	BATTERS_URL  = "http://www.baseballthinkfactory.org/szymborski/ZiPS2012v1BAT.csv"
	BATTERS_CSV  = "zips2012batters.csv"
	PITCHERS_URL = "http://www.baseballthinkfactory.org/szymborski/ZiPS2012v1PIT.csv"
	PITCHERS_CSV = "zips2012pitchers.csv"
)

type ZipsClient struct {
	battingStats  *map[PlayerID]StatLine
	pitchingStats *map[PlayerID]StatLine
}

func (zc *ZipsClient) GetStat(player PlayerID, stat StatID) Stat {
	return zc.GetStatLine(player)[stat]
}

func (zc *ZipsClient) GetStatLine(player PlayerID) StatLine {
	statline, ok := (*zc.battingStats)[player]
	if ok {
		return statline
	}

	return (*zc.pitchingStats)[player]
}

func NewZipsClient() (*ZipsClient, error) {
	EnsureCache(BATTERS_URL, BATTERS_CSV)
	EnsureCache(PITCHERS_URL, PITCHERS_CSV)

	battingStats, err := indexBattingStats()
	if err != nil {
		return nil, err
	}

	pitchingStats, err := indexPitchingStats()
	if err != nil {
		return nil, err
	}

	return &ZipsClient{battingStats: battingStats, pitchingStats: pitchingStats}, nil
}

func mapColumnNameToBattingStat() map[ColName]StatID {
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

func mapColumnNameToPitchingStat() map[ColName]StatID {
	return map[ColName]StatID{
		"ER":  P_EARNED_RUNS,
		"ERA": P_EARNED_RUN_AVERAGE,
		"G":   P_GAMES,
		"H":   P_HITS,
		"HR":  P_HOME_RUNS,
		"IP":  P_INNINGS,
		"L":   P_LOSSES,
		"R":   P_RUNS,
		"GS":  P_STARTS,
		"SO":  P_STRIKE_OUTS,
		"BB":  P_WALKS,
		"W":   P_WINS,
	}
}

func mapBattingStatToColumnIndex(
	colNames []string, columnNameToStat map[ColName]StatID) map[StatID]ColIndex {
	statToColumnIndex := map[StatID]ColIndex{}
	for i := range colNames {
		colName := ColName(colNames[i])
		if statId, ok := columnNameToStat[colName]; ok {
			statToColumnIndex[statId] = ColIndex(i)
		}
	}

	return statToColumnIndex
}

func indexBattingStats() (*map[PlayerID]StatLine, error) {
	return indexStats(BATTERS_CSV, mapColumnNameToBattingStat())
}

func indexPitchingStats() (*map[PlayerID]StatLine, error) {
	return indexStats(PITCHERS_CSV, mapColumnNameToPitchingStat())
}

func indexStats(filename string, columnNameToStat map[ColName]StatID) (*map[PlayerID]StatLine, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrailingComma = true // Ok to end in trailing comma
	recs, err := r.ReadAll()

	if err != nil {
		return nil, err
	}

	statToColumnIndex := mapBattingStatToColumnIndex(recs[0], columnNameToStat)

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
