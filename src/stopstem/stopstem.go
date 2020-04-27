package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var stopWordsBool map[string]bool

var stopwordfile = "../../assets/stopwords.txt"
var readfile = "../crawler/spider_result.txt"
var writefile = "../crawler/spider_result_stemmed.txt"

//InputStopWords function to put all the stopwords listed in the .txt file (duplicates are removed)
func InputStopWords() {
	absPath, _ := filepath.Abs(stopwordfile)
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
	return stopWordsBool[s]
}

//StemThemAll function to remove the unneccessary words
func StemThemAll() {
	//read spider_result.txt
	ReadPath, _ := filepath.Abs(readfile)
	crawled, err := ioutil.ReadFile(ReadPath)
	if err != nil {
		panic(err)
	}
	//result := crawled
	var result []byte
	txtlines := bytes.Split(crawled, []byte("\n"))

	result = append(result, txtlines[1]...)

	for _, lines := range txtlines {
		txtwords := bytes.Split(lines, []byte(" "))
		for _, words := range txtwords {
			//if not an url
			if !strings.HasPrefix(string(words), "http") {
				if CheckStopWords(string(words)) {
					continue
				} else {
					words = []byte(string(words) + " ")
					result = append(result, words...)
				}
				//if an url
			} else {
				result = append(result, words...)
			}
		}
		//result = append(result, '\n')
	}

	//write to spider_result_stemmed.txt
	WritePath, _ := filepath.Abs(writefile)
	err = ioutil.WriteFile(WritePath, result, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	InputStopWords()
	CheckStopWords("tralalalala")
	CheckStopWords("above")
	CheckStopWords("trililii")
	CheckStopWords("becomes")
	CheckStopWords("a")
	StemThemAll()
}
