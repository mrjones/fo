package main

import (
	"folib"

	"bufio"
	"flag"
	"fmt"
	"log"
	"strings"
	"os"
)

func Usage() {
	fmt.Println("Usage:")
	fmt.Print("go run fo.go")
	fmt.Print("  --consumerkey <consumerkey>")
	fmt.Println("  --consumersecret <consumersecret>")
	fmt.Println("")
	fmt.Println("In order to get your consumerkey and consumersecret, you must register an 'app' at twitter.com:")
	fmt.Println("https://dev.twitter.com/apps/new")
}

func loadYahooClientOrDie(key, secret, tokenFile string) *folib.YahooClient {
	if len(key) == 0 || len(secret) == 0 {
		fmt.Println("You must set the --consumerkey and --consumersecret flags.")
		fmt.Println("---")
		Usage()
		os.Exit(1)
	}

	return folib.NewYahooClient(key, secret, tokenFile)
}

func main() {
	var consumerKey *string = flag.String(
		"consumerkey",
		"",
		"Consumer Key from Yahoo. See: https://developer.apps.yahoo.com/dashboard/createKey.html")

	var consumerSecret *string = flag.String(
		"consumersecret",
		"",
		"Consumer Secret from Yahoo. See: https://developer.apps.yahoo.com/dashboard/createKey.html")

	var tokenFile *string = flag.String(
		"tokenfile",
		"",
		"A file to stash the auth token")

	var action *string = flag.String(
		"action",
		"optimize",
		"What to do")

	flag.Parse()

	if *action == "optimize" {
		zipsclient, err := folib.NewZipsClient()
		if err != nil {
			log.Fatal(err)
		}

		yahooclient := loadYahooClientOrDie(*consumerKey, *consumerSecret, *tokenFile)
		fo := folib.NewFO(yahooclient, zipsclient)
		fo.Optimize()
	} else if *action == "summarize" {
		yahooclient := loadYahooClientOrDie(*consumerKey, *consumerSecret, *tokenFile)

		games, err := yahooclient.GetGames()
		if err != nil {
			log.Fatal(err)
		}

		for i, game := range(games) {
			fmt.Printf("%d. %s %s\n", i, game.Name, game.Season)
		}

		gameidx := 0
		if len(games) > 1 {
			log.Fatal("I should ask you which game here and set gameidx...")
		}

		leagues, err := yahooclient.GetLeagues(games[gameidx].GameKey)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d leagues\n", len(leagues))
		for i, league := range(leagues) {
			fmt.Printf("%d. %s\n", i, league.Name)
		}

		leagueidx := 0
		if len(leagues) > 1 {
			log.Fatal("I should ask you which game here and set leagueidx...")
		}

		teams, err := yahooclient.GetTeams(leagues[leagueidx].LeagueKey)
		if err != nil {
			log.Fatal(err)
		}

		for i, team := range(teams) {
			icon := ""
			if team.IsMyTeam == 1 {
				icon = "*"
			}
			fmt.Printf("%d. %s %s\n", i, team.Name, icon)
		}

	} else if *action == "interactive" {
		yahooclient := loadYahooClientOrDie(*consumerKey, *consumerSecret, *tokenFile)

		for {
			fmt.Print("> ");
			input := ""

			bio := bufio.NewReader(os.Stdin)

			done := false
			for !done {
				buf, hasMoreInLine, err := bio.ReadLine()
				input = input + string(buf)

				if err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
					done = true
				}

				if !hasMoreInLine {
					done = true
				}
			}

			fmt.Printf("input is '%s'\n", input)
			inputParts := strings.Split(input, " ")
			fmt.Printf("command is '%s'\n", inputParts[0])

			switch inputParts[0] {
			case "quit", "exit":
				os.Exit(0)
			case "yurl":
				if len(inputParts) < 2 {
					fmt.Println("usage: yurl <url>")
					break
				}
				fmt.Printf("Fetching yahoo url '%s'\n", inputParts[1])
				resp, err := yahooclient.Get(inputParts[1])
				if err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
				} else {
					fmt.Printf("%s\n", resp)
				}
			}
		}
	} else if *action == "fg" {
		_, err := folib.NewFanGraphsClient()
		if err != nil {
			log.Fatal(err)
		}

	}
}
