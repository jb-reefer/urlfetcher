package main

import (
	"strings"
	"testing"
)

func TestSameCaseMatching(t *testing.T) {
	testCorpus := "Some text to search"
	searchText := "text"

	start, _ := searcher.IndexString(testCorpus, searchText)

	if start == -1 {
		t.Errorf("Expected to find %v in %v, instead found start index %v", searchText, testCorpus, start)
	}
}

func TestIgnoreCaseMatching(t *testing.T) {
	testCorpus := "Some text to search"
	searchText := "TEXT"

	start, _ := searcher.IndexString(testCorpus, searchText)

	if start == -1 {
		t.Errorf("Expected to find %v in %v, instead found start index %v", searchText, testCorpus, start)
	}
}

func TestIgnoreFunkyCaseMatching(t *testing.T) {
	testCorpus := "Some text to search"
	searchText := "TeXt"

	start, _ := searcher.IndexString(testCorpus, searchText)

	if start == -1 {
		t.Errorf("Expected to find %v in %v, instead found start index %v", searchText, testCorpus, start)
	}
}

func TestPassingUrl(t *testing.T) {

	response, err := getBody("google.com")

	// Make sure a known good URL passes
	if err != nil {
		t.Error(err)
	}

	defer response.Body.Close()
}

func TestGoodResultString(t *testing.T) {

	result := findResult{true, "http://www.google.com"}
	resultString := result.String()

	if !strings.Contains("true", resultString) &&
		strings.Contains("http://www.google.com", resultString) {

		t.Errorf("Recieved incorrect string representation: %v", resultString)

	}
}

func TestFoundTerm(t *testing.T) {
	query := "google"
	url := "google.com"
	result := searchURLForTerm(url, query)

	if !result.FoundMatch {
		t.Errorf("Expected to find %v at url %v", query, url)
	}
}

func TestNotFoundTerm(t *testing.T) {
	query := "google"
	url := "example.com"
	result := searchURLForTerm(url, query)

	if result.FoundMatch {
		t.Errorf("Didn't expect to find %v at url %v", query, url)
	}
}
