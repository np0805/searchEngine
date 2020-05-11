package pagerank

import (
	"math"
	"strings"

	"../crawler"
	"../database"
	"../stopstem"
)

// CalculatePageRank given a damping factor and map of pages, calculate the ranks recursively
func CalculatePageRank(d float64, pages *map[string]*crawler.Page) {
	for i := 0; i < 1000; i++ {
		// converge := make([]bool, 0)
		// for i := 0; i < len(*pages); i++ {
		// 	converge[i] = false
		// }
		for _, page := range *pages {
			var myRank float64 = 1 - d // value for page rank
			var runningSum float64 = 0 // running sum for probablity from its parents
			for _, p := range page.GetParentURL() {
				parentPage, ok := (*pages)[p]
				if ok {
					var parentPR float64 = parentPage.GetPageRank()
					parentTotalChild := float64(len(parentPage.GetChildrenURL()))
					runningSum += (parentPR / parentTotalChild)
				}
			}
			myRank = myRank + d*runningSum
			difference := myRank - page.GetPageRank()
			if math.Abs(difference) < 0.00000000000005 { // showing signs of converging
				break
			}
			page.SetRank(myRank)
		}
	}
}

// CosSim compute the cosine similarity report for a particular document
func CosSim(queryLength, tfidf, docLength float64) float64 {
	return tfidf / (docLength * queryLength)
}

// TitleMatch cek if a given word match the title
func TitleMatch(word []string, pageID int64) (ok bool, queryScore float64) {
	ok = false
	queryScore = 0.0

	for _, w := range word {
		wordScore := 0.0
		title := database.GetTitle(pageID)
		splitTitle := strings.Split(title, " ")
		titleSlice := make([]string, 0)
		for _, q := range splitTitle {
			titleSlice = append(titleSlice, q)
		}
		titleStem := stopstem.StemString(titleSlice)
		for _, t := range titleStem {
			if w == t {
				ok = true
				wordScore++
				// fmt.Println(w, "match in ", pageID)
			}
		}
		queryScore += wordScore
	}
	return ok, queryScore
}
