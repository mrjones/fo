package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type FangraphsClient struct {
}

func NewFangraphsClient() *FangraphsClient {
	return &FangraphsClient{}
}

func datasets() map[string]string {
	return map[string]string{
		"2012zips.bad.data": "http://www.fangraphs.com/projections.aspx?pos=all&stats=bat&type=zips&team=0&players=0",
	}
}

func (f *FangraphsClient) FetchAllData() {
	files := datasets()

	for filename, url := range files {
		get(url, filename)
	}
}

//
// Factor this out into an httpcache class?
//

func get(url, filename string) (string, error) {
	age := fileAge(filename)

	if age == nil {
		log.Printf("Can't find %s. Downloading", filename)
		httpFetchToFile(url, filename)
	} else {
		ageHours := *age / time.Hour
		if ageHours > 24*30 {
			log.Printf("%s is too old (%d h). Downloading.", filename, ageHours)
			httpFetchToFile(url, filename)
		} else {
			log.Printf("%s is fresh enough (%d h). Not downloading.", filename, ageHours)
		}
	}

	bits, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", err
	}

	return string(bits), nil
}

func fileAge(filename string) *time.Duration {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return nil
	}

	age := time.Since(fileinfo.ModTime())

	return &age
}

func httpFetchToFile(url, filename string) error {
	log.Printf("Fetching %s to %s", url, filename)
	body, err := httpGetBody(url)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, []byte(body), 0644)
}

func httpGetBody(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		return "", err
	}

	bits, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(bits), nil
}
