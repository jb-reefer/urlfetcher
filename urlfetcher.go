package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

var (
	searcher = search.New(language.AmericanEnglish, search.IgnoreCase)
	wait     = sync.WaitGroup{}
)

type findResult struct {
	FoundMatch bool
	URL        string
}

func (result findResult) String() string {
	return fmt.Sprintf("Found: %t\tURL: %v", result.FoundMatch, result.URL)
}

func main() {

	query := flag.String("query", "test", "Search word/phrase in webpages")
	urlFilename := flag.String("file", "urls.txt", "CSV formatted file containg urls")
	threads := flag.Int("threads", 20, "Number of threads to use")

	// Load and parse file
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

	// Create workers and process workload
	urlQueue := make(chan string)
	results := make(chan findResult, len(fileData))

	// Spawn workers
	for i := 0; i < *threads; i++ {
		go searchWorker(urlQueue, results, *query)
	}

	// Pump data into queue
	for _, record := range fileData[1:] {
		urlQueue <- record[1]
	}

	// Await and read out results
	wait.Wait()
	close(urlQueue)
	close(results)

	var buffer bytes.Buffer

	fmt.Println("Reading results...")

	for result := range results {
		buffer.WriteString(result.String())
		buffer.WriteString("\n")
	}

	fmt.Println("Read all results.")

	fileErr := ioutil.WriteFile("results.txt", buffer.Bytes(), 0644)

	if fileErr != nil {
		fmt.Println("Could not write out results.txt.")
		fmt.Println(fileErr)
		return
	}

	fmt.Println("All URLs parsed, exiting.")

}

func searchWorker(input chan string, output chan<- findResult, query string) {

	// Grab urls off of the queue
	for url := range input {
		wait.Add(1)
		output <- searchURLForTerm(url, query)
		wait.Done()
	}
}

func getBody(url string) (*http.Response, error) {

	var address string

	// Input strings aren't always dialable, so try a few different forms
	if url[0:3] != "http" {
		address = "http://" + url
	}

	response, requestError := http.Get(address)

	if requestError != nil {

		address = "http://www." + url

		response, requestError = http.Get(address)
	}

	return response, requestError
}

func searchURLForTerm(url string, query string) findResult {

	response, err := getBody(url)

	if err != nil {
		fmt.Printf("Could not resolve %v:\t%v\n", url, err)
		return findResult{false, url}
	}

	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		fmt.Printf("Could not read body for page %v, got error:%v\n", url, readErr)
		return findResult{false, url}
	}

	start, _ := searcher.Index(body, []byte(query))
	return findResult{start != -1, url}
}
