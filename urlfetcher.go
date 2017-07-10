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
	term    = flag.String("query", "test", "Search word/phrase in webpages")
	matcher = search.New(language.AmericanEnglish, search.IgnoreCase)
	wait    = sync.WaitGroup{}
)

// FindResult Results object from worker threads
type FindResult struct {
	FoundMatch bool
	URL        string
}

func (result FindResult) String() string {
	return fmt.Sprintf("Found: %t\tURL: %v\n", result.FoundMatch, result.URL)
}

func main() {

	// Load and parse file
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

	// Create workers and process workload
	urlQueue := make(chan string)
	results := make(chan FindResult)

	// Spawn workers
	for i := 0; i < *threads; i++ {
		go searchWorker(urlQueue, results)
	}

	// Pump data into queue
	for _, record := range fileData[1:] {
		urlQueue <- record[1]
	}

	// Await and read out results
	wait.Wait()
	close(urlQueue)

	var buffer bytes.Buffer

	for result := range results {
		buffer.WriteString(result.String())
		buffer.WriteString("\n")

	}

	close(results)

	ioutil.WriteFile("results.txt", buffer.Bytes(), 0644)

	fmt.Println("All URLs parsed, exiting.")

}

func searchWorker(input chan string, output chan FindResult) {

	// Grab urls off of the queue
	for url := range input {
		wait.Add(1)
		output <- searchURLForTerm(url)
		wait.Done()
	}
}

func searchURLForTerm(url string) FindResult {

	var address string

	if url[0:3] != "http" {
		address = "http://" + url
	}

	response, httpError := http.Get(address)

	// http.Get doesn't detect proto automatically, so try tls and non-tls
	if httpError != nil {

		address = "http://www." + url

		var wwwError error
		response, wwwError = http.Get(address)

		if wwwError != nil {
			fmt.Printf("Could not resolve %v:\t%v\n", url, wwwError)
			return FindResult{false, address}
		}
	}

	defer response.Body.Close()

	// TODO: this error handling
	body, _ := ioutil.ReadAll(response.Body)
	start, _ := matcher.Index(body, []byte(*term))

	result := FindResult{start != -1, address}
	fmt.Print(result)
	return result
}
