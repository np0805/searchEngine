package database

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"../stopstem"

	bolt "go.etcd.io/bbolt"
)

// N total number of documents in the database
const N = 540

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

// TitleMatch cek if a given word match the title
func TitleMatch(word []string, pageID int64) (ok bool, score float64) {
	ok = false
	score = 0.0

	for _, w := range word {
		title := GetTitle(pageID)
		splitTitle := strings.Split(title, " ")
		titleSlice := make([]string, 0)
		for _, q := range splitTitle {
			titleSlice = append(titleSlice, q)
		}
		titleStem := stopstem.StemString(titleSlice)
		for _, t := range titleStem {
			if w == t {
				ok = true
				score++
				// fmt.Println(w, "match in ", pageID)
			}
		}

	}
	return ok, score
}

func PrintTest() {
	pageInfo.View(func(tx *bolt.Tx) error {
		pageInfoBucket := tx.Bucket([]byte(pageInfoBuck))
		c := pageInfoBucket.Cursor()
		fmt.Println("pageInfoBucket")
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				coba := ByteToString(v)
				fmt.Println("key: ", k, "value: ", coba)
			}
		}
		return nil
	})
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

//idf calculate inverse document frequency of a term
func idf(df int, N float64) float64 {
	return math.Log2(N / float64(df))
}

// DocFreqTerm get document frequency of 1 term j
func DocFreqTerm(word string) map[int64]float64 {
	wordFreqMap := make(map[int64]float64)
	fmt.Println("words ", word)
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
