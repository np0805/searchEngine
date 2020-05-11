package database

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// GetTitle get a title given a page id
func GetTitle(pageID int64) string {
	var title string
	pageInfo.View(func(tx *bolt.Tx) error {
		pageInfoBucket := tx.Bucket([]byte(pageInfoBuck))
		value := pageInfoBucket.Get(IntToByte(pageID))
		stringvalue := ByteToString(value)
		title = stringvalue[0]
		return nil
	})
	return title
}

// PrintTest test aja
func PrintTest() {
	err := wordDb.View(func(tx *bolt.Tx) error {

		wordFreqBucket := tx.Bucket([]byte(wordFreqBuck))
		c := wordFreqBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Println("key: ", GetWord(ByteToInt(k)), "value: ", ByteToString(v)[0])
			break
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// GetLinkRank get the computed link-based page rank of a given pageid
func GetLinkRank(pageID int64) (rank float64) {
	err := pageInfo.View(func(tx *bolt.Tx) error {
		pageRankBucket := tx.Bucket([]byte(pageRankBuck))
		value := pageRankBucket.Get(IntToByte(pageID))
		rank = ByteToFloat64(value)
		return nil
	})
	// fmt.Println("err ", err)
	if err != nil {
		return 0.0
	}
	return rank

}

//idf calculate inverse document frequency of a term
func idf(df int, N float64) float64 {
	return math.Log2(N / float64(df))
}

// DocLength calculate the page frequency length, to be used for cosine similarity
func DocLength(pageID int64) float64 {
	keywords := GetPageKeyFreq(pageID)
	N := float64(GetPageNumber())
	docLength := 0.0
	err := wordDb.View(func(tx *bolt.Tx) error {
		wordFreqBucket := tx.Bucket([]byte(wordFreqBuck))
		for _, wordfreq := range keywords {
			ret := strings.Split(wordfreq, " ")
			intWordID, _ := (strconv.Atoi(ret[0]))

			wordID := int64(intWordID)
			intfreq, _ := (strconv.Atoi(ret[1]))
			tf := float64(intfreq)

			// fmt.Println(IntToByte(w))
			value := wordFreqBucket.Get(IntToByte(wordID))
			stringValue := ByteToString(value)

			df := len(stringValue)
			idf := idf(df, N)

			tfidf := tf * idf
			docLength += tfidf * tfidf
			// fmt.Println("idf of ", GetWord(wordID), " is ", idf, "tf ", tf, "docLength ", docLength)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return docLength
}

// DocFreqTerm get document frequency of 1 term j
func DocFreqTerm(word string) map[int64]float64 {
	wordFreqMap := make(map[int64]float64)
	fmt.Println("words ", word)
	N := float64(GetPageNumber())
	// fmt.Println(GetWordId("comput"))
	err := wordDb.View(func(tx *bolt.Tx) error {
		wordFreqBucket := tx.Bucket([]byte(wordFreqBuck))

		// for _, word := range words {

		w := GetWordId(word)
		if w != 0 { // not found
			fmt.Println(IntToByte(w))
			value := wordFreqBucket.Get(IntToByte(w))
			stringValue := ByteToString(value)
			fmt.Println(stringValue)
			fmt.Println(len(stringValue))

			df := len(stringValue)
			idf := idf(df, N)
			fmt.Println("idf of ", word, " is ", idf)

			for i := 0; i < len(stringValue); i++ {
				res := strings.Split(stringValue[i], " ")
				p, _ := (strconv.Atoi(res[0]))
				pageID := int64(p)
				f, _ := strconv.Atoi(res[1]) // get the frequency
				freq := float64(f)
				tfidf := freq * idf
				frequency, ok := wordFreqMap[pageID]
				if ok {
					frequency += tfidf
				} else {
					frequency = tfidf
				}
				wordFreqMap[pageID] = frequency
			}
		} else {
			fmt.Println("Word not found")
		}
		// }

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return wordFreqMap
}

// WordToWeightMap return a map with key pageID and
// value tf*idf of terms from the given slice
func WordToWeightMap(words []string) map[int64]float64 {
	wordWeightMap := make(map[int64]float64)
	fmt.Println("words ", words)
	N := float64(GetPageNumber())
	// fmt.Println(GetWordId("comput"))
	err := wordDb.View(func(tx *bolt.Tx) error {
		wordFreqBucket := tx.Bucket([]byte(wordFreqBuck))

		for _, word := range words {

			w := GetWordId(word)
			if w == 0 { // not found
				fmt.Println("Word not found")
				continue
			}
			// fmt.Println(IntToByte(w))
			value := wordFreqBucket.Get(IntToByte(w))
			stringValue := ByteToString(value)
			// fmt.Println(stringValue)
			// fmt.Println(len(stringValue))
			df := len(stringValue)
			idf := idf(df, N)
			fmt.Println("idf of ", word, " is ", idf)
			for i := 0; i < len(stringValue); i++ {
				res := strings.Split(stringValue[i], " ")
				p, _ := (strconv.Atoi(res[0])) // get the page id
				pageID := int64(p)
				f, _ := strconv.Atoi(res[1]) // get the frequency
				freq := float64(f)
				tfidf := freq * idf
				frequency, ok := wordWeightMap[pageID]
				if ok {
					frequency += tfidf
				} else {
					frequency = tfidf
				}
				wordWeightMap[pageID] = frequency
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return wordWeightMap
}
