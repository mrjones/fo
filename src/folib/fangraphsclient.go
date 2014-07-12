package folib

import (
	"fmt"
	"time"
)

const (
	ONE_DAY = 24 * time.Hour
)

func fgCsvUrl(year int, tabSet int, batOrPit string) string {
	return fmt.Sprintf("http://www.fangraphs.com/leaders.aspx?pos=all&stats=%s&lg=all&qual=y&type=%d&season=%d&month=0&season1=%d&ind=0&team=0&rost=0&age=0&filter=&players=0", batOrPit, tabSet, year, year)
}

type FanGraphsClient struct {
	battingStats  *map[PlayerID]StatLine
	pitchingStats *map[PlayerID]StatLine
}

func (fc *FanGraphsClient) GetStat(player PlayerID, stat StatID) Stat {
	return fc.GetStatLine(player)[stat]
}

func (fc *FanGraphsClient) GetStatLine(player PlayerID) StatLine {
	statline, ok := (*fc.battingStats)[player]
	if ok {
		return statline
	}

	return (*fc.pitchingStats)[player]
}

func NewFanGraphsClient() (*FanGraphsClient, error) {
	cache := NewReadThroughCache(NewFileKVStore("./cache"))
	_, err := cache.GetAsReader(
		urlFetcher(fgCsvUrl(2014, 0, "bat")), "fg.csv", ONE_DAY)

	if err != nil {
		return nil, err
	}


	return nil, fmt.Errorf("not done")
}


/*
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

func urlFetcher(url string) FetchFunction {
	return func() (string, error) {
		log.Printf("Fetching URL: '%s'", url)
		return httpGetBody(url)
	}
}

func indexBattingStats() (*map[PlayerID]StatLine, error) {
	cache := NewReadThroughCache(NewFileKVStore("./cache"))
	cacheReader, err := cache.GetAsReader(
		urlFetcher(BATTERS_URL), BATTERS_CSV, ONE_MONTH)

	if err != nil {
		return nil, err
	}

	return indexStats(cacheReader, mapColumnNameToBattingStat())
}

func indexPitchingStats() (*map[PlayerID]StatLine, error) {
	cache := NewReadThroughCache(NewFileKVStore("./cache"))
	cacheReader, err := cache.GetAsReader(
		urlFetcher(PITCHERS_URL), PITCHERS_CSV, ONE_MONTH)

	if err != nil {
		return nil, err
	}

	return indexStats(cacheReader, mapColumnNameToPitchingStat())
}

func indexStats(f io.Reader, columnNameToStat map[ColName]StatID) (*map[PlayerID]StatLine, error) {
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
*/
