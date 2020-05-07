package pagerank

import (
	"../crawler"
)

// CalculatePageRank given a damping factor and map of pages, calculate the ranks recursively
func CalculatePageRank(d float64, pages *map[string]*crawler.Page) {
	for i := 0; i < 10; i++ {
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
			page.SetRank(myRank)
		}
	}
}
