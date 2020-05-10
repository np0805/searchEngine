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
	wordMap := database.WordToFreqMap(sli)
	fmt.Println("----")
	fmt.Println(wordMap)
	// for k, page := range wordMap {
	// 	fmt.Println(k, page)
	// }
}
