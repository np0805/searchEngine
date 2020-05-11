package main

import (
	"fmt"

	"./database"
	"./retrieval"
)

func main() {
	database.OpenAllDb()
	r := retrieval.RetrievalFunction("Lip Reading function cibay")
	fmt.Println(r)
	// database.PrintTest()
	// title := database.TitleMatch(510)

	// fmt.Println(title)
	// page := database.GetPageId("https://www.cse.ust.hk/")
	// fmt.Println(page)
	// slice := []string{"competition", "automation", "computing"}
	// sli := stopstem.StemString(slice)
	// fmt.Println(database.GetListOfWordId(sli))
	// // can := []int64{2, 4, 5}
	// // fmt.Println(can)
	// // database.PrintWordDb()

	// wordMap := database.WordToWeightMap(sli)
	// fmt.Println("----")
	// if len(wordMap) == 0 {
	// 	fmt.Println("tewas")
	// }
	// fmt.Println(wordMap)
	// for k, v := range wordMap {
	// 	fmt.Println("key ", k)
	// 	fmt.Println("value ", v)
	// 	break
	// }

	// database.PrintTest()
	// str := database.FindParent("https://www.cse.ust.hk/event/BDICareerFair2020/")
	// fmt.Println(str)
	// for k, page := range wordMap {
	// 	fmt.Println(k, page)
	// }
}
