package database

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// InnerProduct calculate the similarity of the page with the query
func InnerProduct(pageID int64, words []string) (sim float64) {
	return 6.0
}

// TermFreq frequency of term j in document i
func TermFreq(wordID int64, pageID int) int {
	idToByte := IntToByte(wordID)
	frequency := 0
	err := wordDb.View(func(tx *bolt.Tx) error {

		wordFreqBucket := tx.Bucket([]byte(wordFreqBuck))
		value := wordFreqBucket.Get(idToByte)
		// fmt.Println("key: ", wordID, "value: ", value)
		for i := 0; i < len(value); i++ {
			if i%2 == 0 {
				// fmt.Println(value[i])
				if int(value[i]) == pageID {
					frequency = int(value[i+1])
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return frequency
}

// docFreqTerm get document frequency of term j
func docFreqTerm(wordID int64) int {
	idToByte := IntToByte(wordID)
	var length int
	err := wordDb.View(func(tx *bolt.Tx) error {

		wordFreqBucket := tx.Bucket([]byte(wordFreqBuck))
		value := wordFreqBucket.Get(idToByte)
		length = len(value) / 2

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("asd", length)
	return length
}

// // InvDocFreq inverse document frequency of term j
// func InvDocFreq(wordID int64) int {
// 	df := docFreqTerm(wordID)
// 	length := GetDbLength()
// 	fmt.Println("www", length)
// 	k := float64(length / df)
// 	idf := math.Log2(k)
// 	return int(idf)
// }

// PageQueryMap create a map with key page id and value word frequency of the terms in query
func PageQueryMap() {
	err := wordDb.View(func(tx *bolt.Tx) error {

		wordFreqBucket := tx.Bucket([]byte(wordFreqBuck))
		c := wordFreqBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Println("key: ", GetWord(ByteToInt(k)), "value: ", len(v))
			// fmt.Println("test ", )
			break
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// WordToFreqMap return a map with key pageID and
// value frequency of terms from the given slice
func WordToFreqMap(words []string) map[int64]int {
	wordFreqMap := make(map[int64]int)
	fmt.Println("words ", words)
	// fmt.Println(GetWordId("comput"))
	err := wordDb.View(func(tx *bolt.Tx) error {
		wordFreqBucket := tx.Bucket([]byte(wordFreqBuck))

		for _, word := range words {

			w := GetWordId(word)
			if w == 0 { // not found
				fmt.Println("Word not found")
				continue
			}
			fmt.Println(IntToByte(w))
			value := wordFreqBucket.Get(IntToByte(w))
			stringValue := ByteToString(value)
			fmt.Println(stringValue)
			fmt.Println(len(stringValue))
			for i := 0; i < len(stringValue); i++ {
				res := strings.Split(stringValue[i], " ")
				p, _ := (strconv.Atoi(res[0]))
				pageID := int64(p)
				freq, _ := strconv.Atoi(res[1])
				frequency, ok := wordFreqMap[pageID]
				if ok {
					frequency += freq
				} else {
					frequency = freq
				}
				wordFreqMap[pageID] = frequency
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return wordFreqMap
}
