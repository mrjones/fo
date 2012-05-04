package main

import (
	"flag"
	"fmt"
	"log"
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

	flag.Parse()

	if len(*consumerKey) == 0 || len(*consumerSecret) == 0 {
		fmt.Println("You must set the --consumerkey and --consumersecret flags.")
		fmt.Println("---")
		Usage()
		os.Exit(1)
	}

	yahooclient := NewYahooClient(*consumerKey, *consumerSecret, *tokenFile)
	zipsclient, err := NewZipsClient()
	if err != nil {
		log.Fatal(err)
	}

	fo := NewFO(yahooclient, zipsclient)
	fo.Optimize()
}
