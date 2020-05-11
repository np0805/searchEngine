package main

import (
	"fmt"

	"./database"
)

func main() {
	database.OpenAllDb()
	page := database.GetPageId("https://www.cse.ust.hk/")
	fmt.Println(page)
	sli := []string{"compet", "interview", "ppq", "cibay"}
	fmt.Println(database.GetListOfWordId(sli))
	// can := []int64{2, 4, 5}
	// fmt.Println(can)
	// database.PrintWordDb()

	// wordMap := database.DocFreqTerm("comet")
	// fmt.Println("----")
	// if len(wordMap) == 0 {
	// 	fmt.Println("tewas")
	// }
	// fmt.Println(wordMap)

	// database.PrintTest()
	str := database.FindParent("https://www.cse.ust.hk/event/BDICareerFair2020/")
	fmt.Println(str)
	// for k, page := range wordMap {
	// 	fmt.Println(k, page)
	// }
}
