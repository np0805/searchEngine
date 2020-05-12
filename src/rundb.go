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
	r := retrieval.RetrievalFunction("computer -covid19 -more")
	fmt.Println(time.Now())
	for k, v := range r {
		fmt.Println("key ", k, "value", v.GetURL(), v.GetTitle(), v.GetPageRank())
		break
	}
	fmt.Println(len(r))
	// r = retrieval.RetrievalFunction("hkust")
	// fmt.Println(time.Now())
	// for k, v := range r {
	// 	fmt.Println("key ", k, "value", v.GetURL(), v.GetTitle(), v.GetPageRank())
	// 	break
	// }
	// fmt.Println(len(r))
}
