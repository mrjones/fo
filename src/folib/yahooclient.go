package folib

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mrjones/oauth"
)

//
// API
//

func NewYahooClient(consumerKey, consumerSecret, tokenFile string) *YahooClient {
	return &YahooClient{
		tokenFile: tokenFile,
		oauth: oauth.NewConsumer(
			consumerKey,
			consumerSecret,
			oauth.ServiceProvider{
				RequestTokenUrl:   "https://api.login.yahoo.com/oauth/v2/get_request_token",
				AuthorizeTokenUrl: "https://api.login.yahoo.com/oauth/v2/request_auth",
				AccessTokenUrl:    "https://api.login.yahoo.com/oauth/v2/get_token",
			}),
		cache: NewReadThroughCache(NewFileKVStore("./cache")),
	}
}

func (yc *YahooClient) Get(url string) (string, error) {
	fmt.Printf("Getting '%s'\n", url)
	token, err := yc.getAccessToken()
	if err != nil {
		return "", err
	}

	response, err := yc.oauth.Get(
		url,
		map[string]string{},
		token)

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(bits), nil
}

func (yc *YahooClient) GetGames() ([]YahooGame, error) {
	body, err := yc.Get("http://fantasysports.yahooapis.com/fantasy/v2/game/mlb")
	if err != nil {
		return []YahooGame{}, err
	}

	var data FantasyContent
	err = xml.Unmarshal([]byte(body), &data)
	if err != nil {
		return []YahooGame{}, err
	}

	return data.Games, nil
}

func (yc *YahooClient) GetLeagues(gameKey string) ([]YahooLeague, error) {
	url := fmt.Sprintf("http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games;game_keys=%s/leagues", gameKey)
	body, err := yc.Get(url)
	if err != nil {
		return []YahooLeague{}, err
	}

	var data FantasyContent
	err = xml.Unmarshal([]byte(body), &data)
	if err != nil {
		return []YahooLeague{}, err
	}

	if len(data.Users) != 1 {
		return []YahooLeague{}, fmt.Errorf("Wrong number of matching users: %d", len(data.Users))
	}

	if len(data.Users[0].Games) != 1 {
		return []YahooLeague{}, fmt.Errorf("Wrong number of matching games: %d", len(data.Users[0].Games))
	}
	
	return data.Users[0].Games[0].Leagues, nil
}

type listTeamsReply struct {
	Teams []YahooTeam `xml:"leagues>league>teams>team"`
}

func (yc* YahooClient) GetTeams(leagueKey string) ([]YahooTeam, error) {
	url := fmt.Sprintf("http://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys=%s/teams", leagueKey)

	body, err := yc.Get(url)
	if err != nil {
		return []YahooTeam{}, err
	}

	var data listTeamsReply
	err = xml.Unmarshal([]byte(body), &data)
	if err != nil {
		return []YahooTeam{}, err
	}

	return data.Teams, nil
}

type getRosterReply struct {
	Players []YahooPlayer `xml:"team>roster>players>player"`
}

func (yc* YahooClient) GetRoster(teamKey string) ([]YahooPlayer, error) {
	url := fmt.Sprintf("http://fantasysports.yahooapis.com/fantasy/v2/team/%s/roster", teamKey)

	body, err := yc.Get(url)
	if err != nil {
		return []YahooPlayer{}, err
	}

	var data getRosterReply
	err = xml.Unmarshal([]byte(body), &data)
	if err != nil {
		return []YahooPlayer{}, err
	}

	return data.Players, nil	
}

type getStatsReply struct {
	Stats []YahooStat `xml:"player>player_stats>stats>stat"`
}

func (yc* YahooClient) GetStats(playerKey string) (StatLine, error) {
	url := fmt.Sprintf("http://fantasysports.yahooapis.com/fantasy/v2/player/%s/stats", playerKey)

	body, err := yc.Get(url)
	if err != nil {
		return StatLine{}, err
	}

	var data getStatsReply
	err = xml.Unmarshal([]byte(body), &data)
	if err != nil {
		return StatLine{}, err
	}

	statline := StatLine{}
	yahooIdToStatIdMap := mapYahooIdToStatId()
	for _, ystat := range(data.Stats) {
		statid, ok := yahooIdToStatIdMap[ystat.ID]
		if ok {
			statval, err := strconv.ParseFloat(ystat.Value, 64)
			if err != nil {
				return nil, err
			}
			statline[statid] = Stat(statval)
		}
	}

	return statline, nil
}

// old api (hardcoded)

func (yc *YahooClient) CurrentStats() (*map[TeamID]StatLine, error) {
	response, err := yc.cacheGet(
		"current_stats",
		"http://fantasysports.yahooapis.com/fantasy/v2/league/308.l.21006/standings")
//		"http://fantasysports.yahooapis.com/fantasy/v2/league/mlb.l.5181/standings")

	if err != nil {
		return nil, err
	}

	//	fmt.Println(response)

	var data FantasyContent
	err = xml.Unmarshal([]byte(response), &data)
	if err != nil {
		return nil, err
	}

	yahooIdToStatIdMap := mapYahooIdToStatId()

	teamstats := map[TeamID]StatLine{}

	for i := range data.League.Teams {
		team := data.League.Teams[i]
		statline := make(StatLine)
		for j := range team.Stats {
			stat := team.Stats[j]
			statid, ok := yahooIdToStatIdMap[stat.ID]
			if ok {
				statval, err := strconv.ParseFloat(stat.Value, 64)
				if err != nil {
					return nil, err
				}
				statline[statid] = Stat(statval)
			}
		}
		teamstats[team.TeamId] = statline
	}
	return &teamstats, nil
}

func (yc *YahooClient) MyStats() (*StatLine, error) {
	leaguestats, err := yc.CurrentStats()
	if err != nil {
		return nil, err
	}
	mystats := (*leaguestats)[6]
	return &mystats, nil
}

func (yc *YahooClient) LeagueRosters() (*map[TeamID][]YahooPlayer, error) {
	response, err := yc.cacheGet(
		"league_rosters",
		"http://fantasysports.yahooapis.com/fantasy/v2/league/308.l.21006/teams/roster")
//		"http://fantasysports.yahooapis.com/fantasy/v2/league/mlb.l.5181/teams/roster")

	if err != nil {
		return nil, err
	}

	//	fmt.Println(response)

	var data FantasyContent
	err = xml.Unmarshal([]byte(response), &data)
	if err != nil {
		return nil, err
	}

	rosters := map[TeamID][]YahooPlayer{}
	for i := range data.League.Teams2 {
		team := data.League.Teams2[i]
		rosters[team.TeamId] = team.Roster
	}

	return &rosters, nil
}

func (yc *YahooClient) MyRoster() (*[]YahooPlayer, error) {
	response, err := yc.cacheGet(
		"my_roster",
		"http://fantasysports.yahooapis.com/fantasy/v2/team/308.l.21006.t.6/roster")
//		"http://fantasysports.yahooapis.com/fantasy/v2/team/mlb.l.5181.t.6/roster")

	if err != nil {
		return nil, err
	}

	//	fmt.Println(response)

	var data FantasyContent
	err = xml.Unmarshal([]byte(response), &data)
	if err != nil {
		return nil, err
	}

	return &data.Team.Roster, nil
}

//
// Structures
//

type YahooClient struct {
	tokenFile string
	oauth     *oauth.Consumer
	cache     ReadThroughCache
}

type FantasyContent struct {
  Team   YahooTeam   `xml:"team"`
	League YahooLeague `xml:"league"`
	Games  []YahooGame `xml:"game"`
	Users  []YahooUser `xml:"users>user"`
}

type YahooUser struct {
  Games []YahooGame `xml:"games>game"`
}

type YahooGame struct {
	GameKey string `xml:"game_key"`
	GameId string `xml:"game_id"`
	Name string `xml:"name"`
	Season string `xml:"season"`

	Leagues []YahooLeague `xml:"leagues>league"`
}

type YahooTeam struct {
	Name    string `xml:"name"`
	TeamKey string `xml:"team_key"`
	TeamId  TeamID `xml:"team_id"`
	IsMyTeam int `xml:"is_owned_by_current_login"`

	Roster []YahooPlayer `xml:"roster>players>player"`

	Stats []YahooStat `xml:"team_stats>stats>stat"`
}

type YahooStat struct {
	ID    int    `xml:"stat_id"`
	Value string `xml:"value"`
}

type YahooPlayer struct {
	PlayerKey    string   `xml:"player_key"`
	FullName     string   `xml:"name>full"`
	PositionType string   `xml:"position_type"`
	Position     []string `xml:"eligible_positions>position"`
}

type YahooLeague struct {
	Teams     []YahooTeam `xml:"standings>teams>team"`
	Teams2    []YahooTeam `xml:"teams>team"` // OMG :(
	LeagueKey string      `xml:"league_key"`
	Id        int         `xml:"league_id"`
	Name      string      `xml:"name"`
}

//
// Implementation
//

func (yc *YahooClient) Try(key, url string) (string, error){
	return yc.cacheGet(key, url)
}

func oauthUrlFetcher(yc *YahooClient, url string) FetchFunction {
	return func() (string, error) {
		log.Printf("Fetching (via OAuth): '%s'", url)
		return yc.Get(url)
	}
}

func (yc *YahooClient) cacheGet(key string, url string) (string, error) {
	return yc.cache.Get(oauthUrlFetcher(yc, url), key, time.Hour*24)
}


func mapYahooIdToStatId() map[int]StatID {
	return map[int]StatID{
		7:  B_RUNS,
		12: B_HOME_RUNS,
		13: B_RUNS_BATTED_IN,
		16: B_STOLEN_BASES,
		3:  B_BATTING_AVG,
		50: P_INNINGS,
		28: P_WINS,
		32: P_SAVES,
		42: P_STRIKE_OUTS,
		26: P_EARNED_RUN_AVERAGE,
		27: P_WHIP,
	}
}


func (yc *YahooClient) getAccessToken() (*oauth.AccessToken, error) {
	savedBytes, err := ioutil.ReadFile(yc.tokenFile)
	savedString := string(savedBytes)
	var accessToken *oauth.AccessToken
	if err == nil && len(savedString) > 0 {
		accessToken, err = accessTokenFromPlainString(savedString)

		if err != nil {
			return nil, err
		}
	} else {
		requestToken, url, err := yc.oauth.GetRequestTokenAndUrl("oob")

		fmt.Println("(1) Go to: " + url)
		fmt.Println("(2) Grant access, you should get back a verification code.")
		fmt.Println("(3) Enter that verification code here: ")

		verificationCode := ""
		fmt.Scanln(&verificationCode)

		accessToken, err = yc.oauth.AuthorizeToken(requestToken, verificationCode)

		if err != nil {
			return nil, err
		}
		ioutil.WriteFile(yc.tokenFile, []byte(toPlainString(*accessToken)), 0644)
	}

	return accessToken, nil
}

func toPlainString(t oauth.AccessToken) string {
	return fmt.Sprintf("%d|%s|%d|%s", len(t.Token), t.Token, len(t.Secret), t.Secret)
}

func accessTokenFromPlainString(s string) (*oauth.AccessToken, error) {
	firstBar := strings.Index(s, "|")
	if firstBar == -1 {
		return nil, fmt.Errorf("Malformed input [%s]. Couldn't find first bar.", s)
	}

	len1, err := strconv.Atoi(s[0:firstBar])
	if err != nil {
		return nil, fmt.Errorf("Malformed input [%s]", s)
	}
	token := s[firstBar+1 : firstBar+1+len1]

	secondBar := firstBar + len1 + 1
	if s[secondBar] != '|' {
		return nil, fmt.Errorf("Malformed input [%s]. Char %d is not '|'", s, secondBar)
	}

	secondHalf := s[secondBar+1:]

	thirdBar := strings.Index(secondHalf, "|")
	if thirdBar == -1 {
		return nil, fmt.Errorf("Malformed input [%s]. Couldn't find third bar.", s)
	}

	len2, err := strconv.Atoi(secondHalf[0:thirdBar])
	if err != nil {
		return nil, fmt.Errorf("Malformed input [%s]", s)
	}
	secret := secondHalf[thirdBar+1 : thirdBar+1+len2]

	return &oauth.AccessToken{Token: token, Secret: secret}, nil
}

// yurl http://fantasysports.yahooapis.com/fantasy/v2/game/mlb
//   <game_key>328</game_key>
//
// yurl http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games;game_keys=328/leagues
//   <league_key>328.l.1305</league_key>
//   <league_id>1305</league_id>
//   <name>Princeton Sucks</name>
//
// yurl http://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys=328.l.1305
// ...
//
// yurl http://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys=328.l.1305/teams
//    <team>
//      <team_key>328.l.1305.t.5</team_key>
//      <team_id>5</team_id>
//      <name>Curse of Andino</name>
//      <is_owned_by_current_login>1</is_owned_by_current_login>
//
// yurl http://fantasysports.yahooapis.com/fantasy/v2/team/328.l.1305.t.5/roster
//   <player>
//     <player_key>328.p.8395</player_key>
//     <player_id>8395</player_id>
//     <name>
//       <full>Matt Wieters</full>
//
// yurl http://fantasysports.yahooapis.com/fantasy/v2/player/328.p.8395/metadata
// yurl http://fantasysports.yahooapis.com/fantasy/v2/player/328.p.8395/stats
// 



// Full Docs:
// http://developer.yahoo.com/fantasysports/guide/index.html
//
// League standings: 
// "http://fantasysports.yahooapis.com/fantasy/v2/league/mlb.l.5181/standings",
//
// Team Roster:
// "http://fantasysports.yahooapis.com/fantasy/v2/team/mlb.l.5181.t.6/roster",
//
// 10 Free Agents:
// "http://fantasysports.yahooapis.com/fantasy/v2/league/mlb.l.5181/players;status=FA;count=10",
