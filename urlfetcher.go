package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

type FindResult struct {
	FoundMatch bool
	URL        string
}

var matcher = search.New(language.AmericanEnglish, search.IgnoreCase)

func main() {

	var urlFilename = flag.String("file", "urls.txt", "CSV formatted file containg urls")
	var term = flag.String("query", "test", "Search word/phrase in webpages")

	file, err := os.Open(*urlFilename)

	if err != nil {
		fmt.Printf("Error when attempting to open file %v:\n", urlFilename)
		fmt.Println(err)
		return
	}

	csvReader := csv.NewReader(file)
	fileData, err := csvReader.ReadAll()

	if err != nil {
		fmt.Printf("Error when attempting to read file %v:\n", urlFilename)
		fmt.Println(err)
		return
	}

	response := make(chan FindResult, 20)

	// TODO: run 20 of these concurrently
	for _, record := range fileData[1:] {

		concurrentSearchForTerm(record[1], term, response)

	}
}

func concurrentSearchForTerm(url string, term *string, channel chan<- FindResult) {

	channel <- searchURLForTerm(url, term)

}

func searchURLForTerm(url string, term *string) FindResult {

	var address string

	if url[0:3] != "http" {
		address = "https://" + url
	}

	response, httpsError := http.Get(address)

	// http.Get doesn't detect proto automatically, so try tls and non-tls
	if httpsError != nil {

		address = "http://" + url

		var httpError error
		response, httpError = http.Get(address)

		if httpError != nil {

			fmt.Print(httpError)

			return FindResult{false, address}
		}

	}

	defer response.Body.Close()

	// TODO: this error handling
	body, _ := ioutil.ReadAll(response.Body)
	start, _ := matcher.Index(body, []byte(*term))

	return FindResult{start != -1, address}
}
