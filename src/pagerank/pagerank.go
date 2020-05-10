package pagerank

import (
	"math"

	"../crawler"
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
