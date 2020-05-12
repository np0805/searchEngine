package main

import (
	"fmt"
	"time"

	"./database"
	"./retrieval"
)

func main() {
	database.OpenAllDb()
	fmt.Println(time.Now())
	r := retrieval.RetrievalFunction("department of + -covid19  2020 science hkust ")
	fmt.Println(time.Now())
	for k, v := range r { // get the top page
		fmt.Println("key ", k, "value", v.GetURL(), v.GetTitle(), v.GetPageRank(), v.GetID())
		break
	}
	fmt.Println("-------")
	// database.PrintTest()
	fmt.Println(len(r))
	// fmt.Println(database.GetPageNumber())

}
