package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mrjones/oauth"
)

type YahooClient struct {
	tokenFile string
	oauth     *oauth.Consumer
}


type FantasyContent struct {
	Team YahooTeam `xml:"team"`
	League YahooLeague `xml:"league"`
}

type YahooTeam struct {
	Name string `xml:"name"`
	TeamKey string `xml:"team_key"`
	TeamId string `xml:"team_id"`

	Roster []YahooPlayer `xml:"roster>players>player"`

	Stats []YahooTeamStats `xml:"team_stats>stats>stat"`
}

type YahooTeamStats struct {
	ID int `xml:"stat_id"`
	Value string `xml:"value"`
}

type YahooPlayer struct {
	Key string `xml:"player_key"`
	FullName string `xml:"name>full"`
	PositionType string `xml:"position_type"`
	Position []string `xml:"eligible_positions>position"`
}

type YahooLeague struct {
	Teams []YahooTeam `xml:"standings>teams>team"`
	LeagueKey string `xml:"league_key"`
	Id int `xml:"league_id"`
}

func (yc *YahooClient) MyRoster() (*[]YahooPlayer,error) {
	response, err := yc.Get(
		"http://fantasysports.yahooapis.com/fantasy/v2/team/mlb.l.5181.t.6/roster")

	if err != nil { return nil, err }

//	fmt.Println(response)

	var data FantasyContent
	err = xml.Unmarshal([]byte(response), &data)
	if err != nil { return nil, err }

	return &data.Team.Roster, nil
}


func (yc *YahooClient) MyStats() (*YahooTeamStats, error) {
	response, err := yc.Get(
		"http://fantasysports.yahooapis.com/fantasy/v2/league/mlb.l.5181/standings")

	if err != nil { return nil, err }

//	fmt.Println(response)

	var data FantasyContent
	err = xml.Unmarshal([]byte(response), &data)
	if err != nil { return nil, err }

	fmt.Printf("%+v \n", data)

	return nil, nil
}


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
	}
}

func (yc *YahooClient) Get(url string) (string, error) {
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

func (yc *YahooClient) getAccessToken() (*oauth.AccessToken, error) {
	savedBytes, err := ioutil.ReadFile(yc.tokenFile)
	savedString := string(savedBytes)
	var accessToken *oauth.AccessToken
	if err == nil && len(savedString) > 0 {
		accessToken, err = AccessTokenFromPlainString(savedString)

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
		ioutil.WriteFile(yc.tokenFile, []byte(ToPlainString(*accessToken)), 0644)
	}

	return accessToken, nil
}

func ToPlainString(t oauth.AccessToken) string {
	return fmt.Sprintf("%d|%s|%d|%s", len(t.Token), t.Token, len(t.Secret), t.Secret)
}

func AccessTokenFromPlainString(s string) (*oauth.AccessToken, error) {
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
