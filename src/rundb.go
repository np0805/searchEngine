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
	r := retrieval.RetrievalFunction("hkust engine computer")
	fmt.Println(time.Now())
	for k, v := range r {
		fmt.Println("key ", k, "value", v.GetURL(), v.GetTitle(), v.GetPageRank())
		break
	}
	fmt.Println(len(r))
	fmt.Println(database.GetPageNumber())
}
