package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var stopWordsBool map[string]bool

var file = "../assets/stopwords.txt"

//InputStopWords function to put all the stopwords listed in the .txt file (duplicates are removed)
func InputStopWords() {
	absPath, _ := filepath.Abs(file)
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	newlineRegex := regexp.MustCompile("\r?\n")
	stopWordsString := newlineRegex.ReplaceAllString(string(data), " ")
	stopWordsArr := strings.Split(stopWordsString, " ")
	stopWordsBool = make(map[string]bool)

	for _, word := range stopWordsArr {
		stopWordsBool[word] = true
	}
}

//CheckStopWords function to check whether input s is in the map of stopwords
func CheckStopWords(s string) bool {
	if stopWordsBool == nil {
		InputStopWords()
	}
	fmt.Println(stopWordsBool[s])
	return stopWordsBool[s]
}

func main() {
	InputStopWords()
	CheckStopWords("tralalalala")
	CheckStopWords("above")
	CheckStopWords("trililii")
	CheckStopWords("becomes")
	CheckStopWords("a")
}
