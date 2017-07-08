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

var (
	matcher = search.New(language.AmericanEnglish, search.IgnoreCase)
	term    = flag.String("query", "test", "Search word/phrase in webpages")
)

func main() {

	var urlFilename = flag.String("file", "urls.txt", "CSV formatted file containg urls")
	var threads = flag.Int("threads", 20, "Number of threads to use")

	flag.Parse()

	file, err := os.Open(*urlFilename)

	if err != nil {
		fmt.Printf("Error when attempting to open file %v:\n", *urlFilename)
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

	output := make(chan string)

	// Spawn workers
	for i := 0; i < *threads; i++ {
		go searchWorker(output)
	}

	// Pump data into queue
	for _, record := range fileData[1:] {
		output <- record[1]
	}

	fmt.Println("All URLs parsed, exiting.")

	close(output)

}

func searchWorker(input chan string) {

	// Grab urls off of the queue
	for url := range input {
		searchURLForTerm(url)
	}
}

func searchURLForTerm(url string) {

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

			address = "http://www." + url

			var wwwError error
			response, wwwError = http.Get(address)

			if wwwError != nil {
				fmt.Printf("Could not resolve %v\n", url)
				return
			}

		}

	}

	defer response.Body.Close()

	// TODO: this error handling
	body, _ := ioutil.ReadAll(response.Body)
	start, _ := matcher.Index(body, []byte(*term))

	//return FindResult{start != -1, address}
	fmt.Printf("Found: %t\tURL: %v\n", start != -1, address)
}
