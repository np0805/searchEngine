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
	r := retrieval.RetrievalFunction("love Department OF science")
	fmt.Println(time.Now())
	for k, v := range r {
		fmt.Println("key ", k, "value", v.GetURL(), v.GetTitle(), v.GetPageRank())
		break
	}
	// fmt.Println("-------")
	// database.PrintTest()
	// fmt.Println(len(r))
	// fmt.Println(database.GetPageNumber())

}
