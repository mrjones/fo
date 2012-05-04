package main

import (
	"encoding/csv"
	"os"
	"strconv"
)


type Stat float64
type StatID int32
type ColIndex int
type ColName string

type StatLine map[StatID]Stat

const (
	B_AT_BATS StatID = 1
	B_BATTING_AVG StatID = 2
	B_CAUGHT_STEALING StatID = 3
	B_DOUBLES StatID = 4
	B_GAMES StatID = 5
	B_HITS StatID = 6
	B_HOME_RUNS StatID = 7
	B_ON_BASE_PCT StatID = 8
	B_PLATE_APPS StatID = 9
	B_RUNS StatID = 10
	B_RUNS_BATTED_IN StatID = 11
	B_SLUGGING StatID = 12
	B_STOLEN_BASES StatID = 13
	B_STRIKE_OUTS StatID = 14
	B_TRIPLES StatID = 15
	B_WALKS StatID = 16
)

//type PlayerID struct {
//	FirstName string
//	LastName string
//}

type PlayerID string


type ZipsClient struct {
	battingStats *map[PlayerID]StatLine
}

func NewZipsClient() (*ZipsClient, error) {
	get("http://www.baseballthinkfactory.org/szymborski/ZiPS2012v1BAT.csv", "zips2012batters.csv")
	get("http://www.baseballthinkfactory.org/szymborski/ZiPS2012v1PIT.csv", "zips2012pitchers.csv")

	battingStats, err := indexBattingStats()
	if err != nil { return nil, err }

	return &ZipsClient{battingStats: battingStats}, nil
}

func mapColumnNameToStat() map[ColName]StatID {
	return map[ColName]StatID{
		"AB": B_AT_BATS,
	  "BA": B_BATTING_AVG,
		"CS": B_CAUGHT_STEALING,
		"2B": B_DOUBLES,
		"G": B_GAMES,
	  "H": B_HITS,
		"HR": B_HOME_RUNS,
		"OBP": B_ON_BASE_PCT,
		"PA": B_PLATE_APPS,
		"R": B_RUNS,
		"RBI": B_RUNS_BATTED_IN,
		"SLG": B_SLUGGING,
		"SB": B_STOLEN_BASES,
		"K": B_STRIKE_OUTS,
		"3B": B_TRIPLES,
		"BB": B_WALKS,
	}
}

func mapStatToColumnIndex(colNames []string) map[StatID]ColIndex {
	statToColumnIndex := map[StatID]ColIndex{}
	columnNameToStat := mapColumnNameToStat()
	for i := range(colNames) {
		colName := ColName(colNames[i])
		if statId, ok := columnNameToStat[colName]; ok {
			statToColumnIndex[statId] = ColIndex(i)
		}
	}

	return statToColumnIndex
}

func indexBattingStats() (*map[PlayerID]StatLine, error) {
	f, err := os.Open("zips2012batters.csv")
	if err != nil { return nil, err }

	r := csv.NewReader(f)
	recs, err := r.ReadAll()

	if err != nil { return nil, err }

	statToColumnIndex := mapStatToColumnIndex(recs[0])

	statIndex := map[PlayerID]StatLine{}
	for i := 1 ; i < len(recs); i++ {
		statLine := make(StatLine)
		for stat, index := range(statToColumnIndex) {
			statStr := recs[i][index]
			stat64, err := strconv.ParseFloat(statStr, 64)
			if err != nil { return nil, err }
			statLine[stat] = Stat(stat64)
		}
		statIndex[PlayerID(recs[i][0])] = statLine
	}

	return &statIndex, nil
}

func (zc *ZipsClient)	GetStat(player PlayerID, stat StatID) Stat {
	return (*zc.battingStats)[player][stat]
}

func (zc *ZipsClient)	GetStatLine(player PlayerID) StatLine {
	return (*zc.battingStats)[player]
}
