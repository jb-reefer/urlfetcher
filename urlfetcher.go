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
	term        = flag.String("query", "test", "Search word/phrase in webpages")
	urlFilename = flag.String("file", "urls.txt", "CSV formatted file containg urls")
	threads     = flag.Int("threads", 20, "Number of threads to use")

	matcher = search.New(language.AmericanEnglish, search.IgnoreCase)
	wait    = sync.WaitGroup{}
)

// FindResult for a given term in webpage at URL
type FindResult struct {
	FoundMatch bool
	URL        string
}

func (result FindResult) String() string {
	return fmt.Sprintf("Found: %t\tURL: %v", result.FoundMatch, result.URL)
}

func main() {

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
	results := make(chan FindResult, len(fileData))

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
	close(results)

	var buffer bytes.Buffer

	fmt.Println("Reading results...")

	i := 1

	for result := range results {
		fmt.Println(result)
		fmt.Println(i)
		i++
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

func searchWorker(input chan string, output chan<- FindResult) {

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
			//fmt.Printf("Could not resolve %v:\t%v\n", url, wwwError)
			return FindResult{false, url}
		}
	}

	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		fmt.Printf("Could not read body for page %v, got error:%v\n", url, readErr)
		return FindResult{false, url}
	}

	start, _ := matcher.Index(body, []byte(*term))
	return FindResult{start != -1, address}
}
