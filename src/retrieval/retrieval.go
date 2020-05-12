package retrieval

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"

	"../database"
	"../pagerank"
	"../stopstem"
)

// PageScore struct
// score
// title
// page url
// last modif date
// size
// top keywords
// []parent link
// []children link
type PageScore struct {
	Id           int64
	Score        float64
	Title        string
	Url          string
	LastModified string
	PageSize     string
	Keywords     []string
	Parents      []string
	Children     []string
}

// GetID return id of page
func (page *PageScore) GetID() int64 {
	return page.Id
}

// GetPageRank return score of page
func (page *PageScore) GetPageRank() float64 {
	return page.Score
}

// GetTitle return id of page
func (page *PageScore) GetTitle() string {
	return page.Title
}

// GetURL return url of page
func (page *PageScore) GetURL() string {
	return page.Url
}

// GetLastModified return modified date of page
func (page *PageScore) GetLastModified() string {
	return page.LastModified
}

// GetSize return size of page
func (page *PageScore) GetSize() string {
	return page.PageSize
}

// GetKeywords return top 5 keywords of page
func (page *PageScore) GetKeywords() []string {
	return page.Keywords
}

// GetParents return the parents of page
func (page *PageScore) GetParents() []string {
	return page.Parents
}

// GetChildren return the children of page
func (page *PageScore) GetChildren() []string {
	return page.Children
}

// RetrievalFunction return a slice of pages sorted by similarity score with the query
func RetrievalFunction(query string) []*PageScore {
	// pageScoreMap := make(map[int64]float64)
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+ ")
	query = reg.ReplaceAllString(string(query), " ")
	pagesScores := make([]*PageScore, 0)
	querySlice := make([]string, 0)
	splitQuery := strings.Split(query, " ")

	minusWord := make([]string, 0) // to get the list of words precedeed by "-", if any
	for _, q := range splitQuery {
		if q != "" {
			if string(q[0]) == "-" {
				minusWord = append(minusWord, q[1:len(q)])
			} else {
				querySlice = append(querySlice, q)
			}
		}
	}

	var excludedPage []int64
	if minusWord != nil { // there is 1 or more minus term
		minusWord := stopstem.StemString(minusWord)
		excludedPage = database.GetPage(minusWord)
	}
	queryLength := math.Sqrt(float64(len(querySlice)))

	queryStem := stopstem.StemString(querySlice)
	wordMap := database.WordToWeightMap(queryStem)
	if len(wordMap) == 0 {
		fmt.Println("No result for search using this query")
		return nil
	}

	if excludedPage != nil {
		for _, id := range excludedPage {
			_, ok := wordMap[id]
			if ok {
				delete(wordMap, id)
			}
		}
	}

	for k, v := range wordMap {
		docLength := math.Sqrt(database.DocLength(k))
		_, titleScore := pagerank.TitleMatch(queryStem, k) // check for a match in the title and give boost in ranking
		cossim := pagerank.CosSim(queryLength, v, docLength)
		linkrank := database.GetLinkRank(k)
		title, url, lastmodified, size := database.ExtractPageInfo(k)

		topWords := database.GetTopWords(k)
		parents := database.FindParentById(k)
		children := database.FindChildById(k)
		score := cossim + titleScore + linkrank
		// if url == "https://www.cse.ust.hk/" {
		// 	fmt.Println("-------------")
		// 	fmt.Println("cosin", cossim)
		// 	fmt.Println("title", titleScore)
		// 	fmt.Println("linkrank", linkrank)
		// 	fmt.Println("doclength di retrieve", docLength)
		// 	fmt.Println("-------------")
		// }
		pageScore := PageScore{
			Id:           k,
			Score:        score,
			Title:        title,
			Url:          url,
			LastModified: lastmodified,
			PageSize:     size,
			Keywords:     topWords,
			Parents:      parents,
			Children:     children}

		pagesScores = append(pagesScores, &pageScore)
		// pageScoreMap[k] = cossim + titleScore + linkrank
	}
	// sort.Sort(pagesScores)
	sort.SliceStable(pagesScores, func(i, j int) bool {
		return pagesScores[i].Score > pagesScores[j].Score
	})

	// write to json, kalo ini gadipake dibuang aja
	file, _ := json.MarshalIndent(pagesScores, "", " ")

	ioutil.WriteFile("search_output.json", file, os.ModePerm)

	return pagesScores
}
