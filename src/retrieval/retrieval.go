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
// score
// title
// page url
// last modif date
// size
// top keywords
// []parent link
// []children link
type PageScore struct {
	id           int64
	score        float64
	title        string
	url          string
	lastModified string
	pageSize     string
	keywords     []string
	parents      []string
	children     []string
}

// GetID return id of page
func (page *PageScore) GetID() int64 {
	return page.id
}

// GetScore return score of page
func (page *PageScore) GetScore() float64 {
	return page.score
}

// GetTitle return id of page
func (page *PageScore) GetTitle() string {
	return page.title
}

// GetURL return url of page
func (page *PageScore) GetURL() string {
	return page.url
}

// GetLastModified return modified date of page
func (page *PageScore) GetLastModified() string {
	return page.lastModified
}

// GetSize return size of page
func (page *PageScore) GetSize() string {
	return page.pageSize
}

// GetKeywords return top 5 keywords of page
func (page *PageScore) GetKeywords() []string {
	return page.keywords
}

// GetParents return the parents of page
func (page *PageScore) GetParents() []string {
	return page.parents
}

// GetChildren return the children of page
func (page *PageScore) GetChildren() []string {
	return page.children
}

// RetrievalFunction return a slice of pages sorted by similarity score with the query
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
		title, url, lastmodified, size := database.ExtractPageInfo(k)

		topWords := database.GetTopWords(k)
		parents := database.FindParentById(k)
		children := database.FindChildById(k)
		pageScore := PageScore{
			id:           k,
			score:        cossim + titleScore + linkrank,
			title:        title,
			url:          url,
			lastModified: lastmodified,
			pageSize:     size,
			keywords:     topWords,
			parents:      parents,
			children:     children}

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
