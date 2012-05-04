package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"

	"github.com/mrjones/oauth"
)

type YahooClient struct {
	tokenFile string
	oauth *oauth.Consumer
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
	if err != nil { return "", err }

	response, err := yc.oauth.Get(
		url,
		map[string]string{},
		token)

	if err != nil { return "", err }
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
	if err != nil { return "", err }

	return string(bits), nil
}

func (yc *YahooClient) getAccessToken() (*oauth.AccessToken, error) {
	savedBytes, err := ioutil.ReadFile(yc.tokenFile)
	savedString := string(savedBytes)
	var accessToken *oauth.AccessToken;
	if (err == nil && len(savedString) > 0) {
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

func ToPlainString(t oauth.AccessToken) (string) {
	return fmt.Sprintf("%d|%s|%d|%s", len(t.Token), t.Token, len(t.Secret), t.Secret)
}

func AccessTokenFromPlainString(s string) (*oauth.AccessToken, error) {
	firstBar := strings.Index(s, "|")
	if (firstBar == -1) { return nil, fmt.Errorf("Malformed input [%s]. Couldn't find first bar.", s) }

	len1, err := strconv.Atoi(s[0:firstBar])
	if (err != nil) { return nil, fmt.Errorf("Malformed input [%s]", s) }
	token := s[firstBar + 1:firstBar + 1 + len1]

	secondBar := firstBar + len1 + 1
	if (s[secondBar] != '|') {
		return nil, fmt.Errorf("Malformed input [%s]. Char %d is not '|'", s, secondBar)
	}

	secondHalf := s[secondBar + 1:]

	thirdBar := strings.Index(secondHalf, "|")
	if (thirdBar == -1) { return nil, fmt.Errorf("Malformed input [%s]. Couldn't find third bar.", s) }

	len2, err := strconv.Atoi(secondHalf[0:thirdBar])
	if (err != nil) { return nil, fmt.Errorf("Malformed input [%s]", s) }
	secret := secondHalf[thirdBar + 1:thirdBar + 1 + len2]

	return &oauth.AccessToken{Token: token, Secret: secret}, nil
}
