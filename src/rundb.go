package main

import (
	"fmt"

	"./database"
)

func main() {
	database.OpenAllDb()
	page := database.GetPageId("https://www.cse.ust.hk/")
	fmt.Println(page)
<<<<<<< HEAD
	sli := []string{"compet", "automat", "cibay"}
=======
	sli := []string{"compet", "interview", "ppq", "cibay"}
>>>>>>> 04196080dc6691bfb1ca09edc51eeefff8c9e9b1
	fmt.Println(database.GetListOfWordId(sli))
	// can := []int64{2, 4, 5}
	// fmt.Println(can)
	// database.PrintWordDb()

<<<<<<< HEAD
	wordMap := database.DocFreqTerm("compet")
	fmt.Println("----")
	if len(wordMap) == 0 {
		fmt.Println("tewas")
	}
	fmt.Println(wordMap[18])

	// database.PrintTest()
	// str := database.FindParent("https://www.cse.ust.hk/event/BDICareerFair2020/")
	// fmt.Println(str)
=======
	// wordMap := database.DocFreqTerm("comet")
	// fmt.Println("----")
	// if len(wordMap) == 0 {
	// 	fmt.Println("tewas")
	// }
	// fmt.Println(wordMap)

	// database.PrintTest()
	str := database.FindParent("https://www.cse.ust.hk/event/BDICareerFair2020/")
	fmt.Println(str)
>>>>>>> 04196080dc6691bfb1ca09edc51eeefff8c9e9b1
	// for k, page := range wordMap {
	// 	fmt.Println(k, page)
	// }
}
