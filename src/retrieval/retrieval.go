package retrieval

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"../database"
	"../pagerank"
	"../stopstem"
)

// PageScore struct
type PageScore struct {
	id    int64
	score float64
}

// GetID return id of page
func (page *PageScore) GetID() int64 {
	return page.id
}

// GetScore return url of page
func (page *PageScore) GetScore() float64 {
	return page.score
}

// RetrievalFunction return a map of page id and similarity score given a query
func RetrievalFunction(query string) []*PageScore {
	// pageScoreMap := make(map[int64]float64)
	pagesScores := make([]*PageScore, 0)
	querySlice := make([]string, 0)
	splitQuery := strings.Split(query, " ")
	for _, q := range splitQuery {
		querySlice = append(querySlice, q)
	}
	queryLength := math.Sqrt(float64(len(querySlice)))

	fmt.Println("length", queryLength)
	queryStem := stopstem.StemString(querySlice)
	wordMap := database.WordToWeightMap(queryStem)
	if len(wordMap) == 0 {
		fmt.Println("No result for search using this query")
		return nil
	}

	// fmt.Println(wordMap)
	// i := 0
	for k, v := range wordMap {
		docLength := math.Sqrt(database.DocLength(k))
		_, titleScore := pagerank.TitleMatch(queryStem, k) // check for a match in the title and give boost in ranking
		cossim := pagerank.CosSim(queryLength, v, docLength)
		linkrank := database.GetLinkRank(k)
		pageScore := PageScore{
			id:    k,
			score: cossim + titleScore + linkrank}

		pagesScores = append(pagesScores, &pageScore)
		// pageScoreMap[k] = cossim + titleScore + linkrank
	}
	// sort.Sort(pagesScores)
	sort.SliceStable(pagesScores, func(i, j int) bool {
		return pagesScores[i].score > pagesScores[j].score
	})
	return pagesScores
}

// FillPage fill the values of the PageScore struct
// title
// page url
// last modif date
// size
// top keywords
// []parent link
// []children link
func (page *PageScore) FillPage() {
	// sortedMap := make(map[int64]float64)
}
