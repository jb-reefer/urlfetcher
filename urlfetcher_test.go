package main

import (
	"net/http"
	"strings"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

func TestSameCaseMatching(t *testing.T) {
	testCorpus := "Some text to search"
	searchText := "text"

	var matcher = search.New(language.AmericanEnglish, search.IgnoreCase)

	start, _ := matcher.IndexString(testCorpus, searchText)

	if start == -1 {
		t.Errorf("Expected to find %v in %v, instead found start index %v", searchText, testCorpus, start)
	}
}

func TestIgnoreCaseMatching(t *testing.T) {
	testCorpus := "Some text to search"
	searchText := "TEXT"

	var matcher = search.New(language.AmericanEnglish, search.IgnoreCase)

	start, _ := matcher.IndexString(testCorpus, searchText)

	if start == -1 {
		t.Errorf("Expected to find %v in %v, instead found start index %v", searchText, testCorpus, start)
	}
}

func TestIgnoreFunkyCaseMatching(t *testing.T) {
	testCorpus := "Some text to search"
	searchText := "TeXt"

	var matcher = search.New(language.AmericanEnglish, search.IgnoreCase)

	start, _ := matcher.IndexString(testCorpus, searchText)

	if start == -1 {
		t.Errorf("Expected to find %v in %v, instead found start index %v", searchText, testCorpus, start)
	}
}

func TestFailingUrl(t *testing.T) {

	response, err := http.Get("http://123-reg.co.uk/")

	// I know this url is failing, but it is legal. Not wholey sure why.
	if err == nil {
		t.Error(err)
	}

	defer response.Body.Close()
}

func TestResultString(t *testing.T) {

	result := FindResult{true, "http://www.google.com"}
	resultString := result.String()

	if !strings.Contains("true", resultString) && strings.Contains("http://www.google.com", resultString) {
		t.Errorf("Recieved incorrect string representation: %v", resultString)

	}
}
