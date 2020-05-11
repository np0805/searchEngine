package retrieval

import (
	"fmt"
	"math"
	"strings"

	"../database"
	"../pagerank"
	"../stopstem"
)

// RetrievalFunction return a map of page id and similarity score given a query
// 							BELOM KELAR
func RetrievalFunction(query string) map[int64]float64 {
	pageScoreMap := make(map[int64]float64)
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
		//BELOM KELAR
		// get the length of keywords through k
		// get the pagerank calculation from the db
		_, titleScore := pagerank.TitleMatch(queryStem, k) // check for a match in the title and give boost in ranking
		cossim := pagerank.CosSim(queryLength, v, 1.0)
		linkrank := database.GetLinkRank(k)

		pageScoreMap[k] = cossim + titleScore + linkrank
	}
	return pageScoreMap
}

// SortMap sort the pagemap by its rank
func SortMap(pageScoreMap *map[int64]float64) {

}
