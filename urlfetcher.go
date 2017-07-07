package urlfetcher

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

func main() {

	var urlFilename = flag.String("file", "urls.txt", "CSV formatted file containg urls")
	var term = flag.String("query", "test", "Search word/phrase in webpages")

	file, err := os.Open(*urlFilename)

	if err != nil {
		fmt.Printf("Error when attempting to open file %v:", urlFilename)
		fmt.Print(err)
		return
	}

	csvReader := csv.NewReader(file)
	fileData, err := csvReader.ReadAll()

	if err != nil {
		fmt.Printf("Error when attempting to read file %v:", urlFilename)
		fmt.Print(err)
		return
	}

	// TODO: run 20 of these concurrently
	for _, record := range fileData {

		searchURLForTerm(record[0], term)

	}
}

func searchURLForTerm(url string, term *string) bool {

	response, webError := http.Get("https://www.example.com")

	if webError != nil {
		fmt.Print(webError)
		return false
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	var matcher = search.New(language.AmericanEnglish, search.IgnoreCase)

	start, _ := matcher.Index(body, []byte(*term))

	return start != -1
}
