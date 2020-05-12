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
func GetTitle(pageID int64) (title []string) {
	pageInfo.View(func(tx *bolt.Tx) error {
		pageTitleStemBucket := tx.Bucket([]byte(pageTitleStemBuck))
		value := pageTitleStemBucket.Get(IntToByte(pageID))
		stringvalue := ByteToString(value)
		title = stringvalue
		return nil
	})
	return title
}

// PrintTest test aja
func PrintTest() {
	pageInfo.View(func(tx *bolt.Tx) error {
		pageInfoBucket := tx.Bucket([]byte(pageInfoBuck))
		c := pageInfoBucket.Cursor()

		fmt.Println("pageInfoBucket")
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				fmt.Println("key: ", ByteToInt(k), "value: ", ByteToString(v))
			}
			break
		}

		fmt.Println("pageTitleStemBucket")
		pageTitleStemBucket := tx.Bucket([]byte(pageTitleStemBuck))
		c = pageTitleStemBucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Println("key: ", ByteToInt(k), "value: ", ByteToString(v))
			break
		}

		return nil
	})
}

// ExtractPageInfo get the info from pageinfo bucket
func ExtractPageInfo(pageID int64) (title, url, lastmodified, size string) {
	pageInfo.View(func(tx *bolt.Tx) error {
		pageInfoBucket := tx.Bucket([]byte(pageInfoBuck))
		value := pageInfoBucket.Get(IntToByte(pageID))
		v := ByteToString(value)
		title = v[0]
		url = v[1]
		lastmodified = v[2]
		size = v[3]

		return nil
	})
	return title, url, lastmodified, size
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
	// if pageID == GetPageId("https://www.cse.ust.hk/") {
	// 	fmt.Println("doclegnth", math.Sqrt(docLength))
	// }
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
	// fmt.Println("words ", words)
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
			// fmt.Println("idf of ", word, " is ", idf)
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

// GetPage return pages that contain the words
func GetPage(words []string) (pageID []int64) {
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
			for i := 0; i < len(stringValue); i++ {
				res := strings.Split(stringValue[i], " ")
				p, _ := (strconv.Atoi(res[0])) // get the page id
				pid := int64(p)
				pageID = append(pageID, pid)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return pageID
}

// FindChildById given a page Id, return the child []string
func FindChildById(Id int64) (ret []string) {
	pageId := IntToByte(Id)
	err := pageInfo.View(func(tx *bolt.Tx) error {
		parentChildBucket := tx.Bucket([]byte(parentChildBuck))
		value := parentChildBucket.Get(pageId)
		if value != nil {
			ret = ByteToString(value)
		} else {
			ret = nil
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return ret
}

// given a page Id, return the parent []string
func FindParentById(Id int64) (ret []string) {
	pageId := IntToByte(Id)
	err := pageInfo.View(func(tx *bolt.Tx) error {
		childParentBucket := tx.Bucket([]byte(childParentBuck))
		value := childParentBucket.Get(pageId)
		if value != nil {
			ret = ByteToString(value)
		} else {
			ret = nil
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return ret
}
