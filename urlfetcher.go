package urlfetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

func main() {

	// TODO: take flags in from the args with /flags package
	// TODO: write function that takes a URL and a term, and searches for a value
	// TODO: write function that searches for a value

	response, webError := http.Get("https://www.example.com")

	if webError != nil {
		fmt.Print(webError)
		return
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(body)

	var matcher = search.New(language.AmericanEnglish, search.IgnoreCase)

	matcher.IndexString("term", "corpus")

}
