package main

import (
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
