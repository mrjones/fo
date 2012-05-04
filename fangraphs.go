package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type FangraphsClient struct {
}

func NewFangraphsClient() *FangraphsClient {
	return &FangraphsClient{}
}

func (f *FangraphsClient) GetZipsProjections() (string, error) {
	filename := "zips.data"
	url := "http://www.fangraphs.com/projections.aspx?pos=all&stats=bat&type=zips&team=0&players=0"
	
	age := fileAge(filename)
	if (age == nil || *age / time.Hour > 24 * 30) {
		httpFetchToFile(url, filename)
	}
	
	bits, err := ioutil.ReadFile(filename)
	
	if err != nil { return "", err }

	return string(bits), nil
}

func fileAge(filename string) *time.Duration {
	file, err := os.Open(filename)
	if err != nil { return nil }
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil { return nil }

	age := time.Since(fileinfo.ModTime())

	return &age
}

func httpFetchToFile(url, filename string) error {
	body, err := httpGetBody(url)
	if err != nil { return err }

	return ioutil.WriteFile(filename, []byte(body), 0644)
}

func httpGetBody(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil { return "", err }

	bits, err := ioutil.ReadAll(response.Body)

	if err != nil { return "", err }

	return string(bits), nil
}
