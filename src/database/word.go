package database

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"../crawler"
	bolt "go.etcd.io/bbolt"
)

var wordDb *bolt.DB
var wordToIdBuck string = "wordIdBuck"
var idToWordBuck string = "idToWordBuck"
var wordFreqBuck string = "wordFreqBuck"
var pageWordFreqBuck string = "pageWordFreqBuck"

func openWordDb() {
	var err error
	wordDb, err = bolt.Open("db"+string(os.PathSeparator)+"word.db", 0700, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = wordDb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(wordToIdBuck))
		if err != nil {
			return fmt.Errorf("word create first bucket error: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(idToWordBuck))
		if err != nil {
			return fmt.Errorf("word create second bucket error: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(wordFreqBuck))
		if err != nil {
			return fmt.Errorf("word create third bucket error: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(pageWordFreqBuck))
		if err != nil {
			return fmt.Errorf("word create fourth bucket error: %s", err)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func closeWordDb() {
	wordDb.Close()
}

// given a word in string, return its wordId in int64, returns 0 if word does not exist
func GetWordId(word string) (ret int64) {
	err := wordDb.View(func(tx *bolt.Tx) error {
		wordIdBucket := tx.Bucket([]byte(wordToIdBuck))
		value := wordIdBucket.Get([]byte(word))
		if value != nil {
			ret = ByteToInt(value)
		} else {
			ret = 0
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return ret
}

// given an id int64, return the word string, returns "" if does not exist
func GetWord(id int64) (word string) {
	word = ""
	idToByte := IntToByte(id)
	err := wordDb.View(func(tx *bolt.Tx) error {
		idToWordBucket := tx.Bucket([]byte(idToWordBuck))
		value := idToWordBucket.Get(idToByte)
		if value != nil {
			word = string(value)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return word
}

// given a list of words, get all their ids
func GetListOfWordId(words []string) (wordIds []int64) {
	for _, word := range words {
		wordIds = append(wordIds, GetWordId(word))
	}
	return wordIds
}

// create the wordId for a word and returns the id
func createWordId(word string) (wordId int64) {
	wordId = GetWordId(word)
	// check first if the word already exists, if it does, simply return the wordId, if it doesn't, handle
	if wordId == 0 && word != "" && word != " " {
		err := wordDb.Update(func(tx *bolt.Tx) error {
			idToWordBucket := tx.Bucket([]byte(idToWordBuck))
			id, _ := idToWordBucket.NextSequence()
			wordId = int64(id)
			err := idToWordBucket.Put(IntToByte(wordId), []byte(word))
			if err != nil {
				fmt.Errorf("Error inserting idToWordBucket, word: ", word, "error: %s", err)
			}

			wordToIdBucket := tx.Bucket([]byte(wordToIdBuck))
			err = wordToIdBucket.Put([]byte(word), IntToByte(wordId))

			// fmt.Println("id: ", id, "word: ", word)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	return wordId
}

// given a pageId int64 and word string, update the frequency table
// update the pageId -> list of words n freq too
func updateFreq(pageId int64, word string) {
	// ignore empty phrases or spaces
	if word != " " && word != "" {
		// check if previous data for the word exists
		exists := false
		existsInPageWord := false
		var value []byte = nil
		var pageWordFreqValue []byte = nil
		termId := createWordId(word)
		err := wordDb.View(func(tx *bolt.Tx) error {
			wordFreqBucket := tx.Bucket([]byte(wordFreqBuck))
			value = wordFreqBucket.Get(IntToByte(termId))
			if value != nil {
				exists = true
			}

			pageWordFreqBucket := tx.Bucket([]byte(pageWordFreqBuck))
			pageWordFreqValue = pageWordFreqBucket.Get(IntToByte(pageId))
			if pageWordFreqValue != nil {
				existsInPageWord = true
			}

			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

		// if entry does not exists, add the first entry
		var toInsert []string
		if exists == true {
			// fmt.Println("exist")
			// handle updating the entry
			// first convert the data in the db to a manipulatable format
			oldStringData := ByteToString(value)
			// fmt.Println("URL: " + GetPageUrl(pageId) + "pageId: " + string(strconv.FormatInt(pageId, 10)) + "word: " + word + " wordId: " + string(strconv.FormatInt(termId, 10)))
			// fmt.Println("oldString: ", oldStringData)
			// check if the entry for this current page exist for this word
			// if they exist, increment the keyFreq, then append to toInsert
			// if not, just append immediately
			entryExist := false
			for _, val := range oldStringData {
				pageStat := strings.Split(val, " ")
				docId, err := strconv.ParseInt(pageStat[0], 10, 64)
				if err != nil {
					fmt.Errorf("Error in wordId, converting docId error: %s", err)
				}

				// check if the pageId entry exist
				if docId == pageId {
					entryExist = true
					keyFreq, err := strconv.ParseInt(pageStat[1], 10, 64)
					if err != nil {
						fmt.Errorf("Error in wordId, converting keyFreq error: %s", err)
					}
					keyFreq = keyFreq + int64(1)
					toInsert = append(toInsert, string(strconv.FormatInt(pageId, 10)+" "+strconv.FormatInt(keyFreq, 10)))
					// toInsert = append(toInsert, string("1 12"))
				} else {
					toInsert = append(toInsert, string(val))
					// toInsert = append(toInsert, string(val))
				}
			}
			if entryExist == false {
				toInsert = append(toInsert, string(strconv.FormatInt(pageId, 10)+" "+strconv.FormatInt(int64(1), 10)))
				// toInsert = append(toInsert, string("12 1"))
			}
		} else {
			toInsert = append(toInsert, string(strconv.FormatInt(pageId, 10)+" "+strconv.FormatInt(int64(1), 10)))
			// toInsert = append(toInsert, string("13 28"))
		}

		var toInsertPage []string
		if existsInPageWord {
			// take the old data and modify accordingly
			oldStringPageData := ByteToString(pageWordFreqValue)
			wordExists := false

			for _, val := range oldStringPageData {
				// wordStats[0] will be wordId wordStats[1] will be frequency
				wordStats := strings.Split(val, " ")
				wordId, err := strconv.ParseInt(wordStats[0], 10, 64)
				if err == nil {
					fmt.Errorf("Error in wordId, converting wordId error: %s", err)
				}

				// check if the word is in the entry of pageWordFreq
				if termId == wordId {
					wordExists = true
					wordFreq, err := strconv.ParseInt(wordStats[1], 10, 64)
					if err == nil {
						fmt.Errorf("Error in wordId, converting wordFreq error: %s", err)
					}
					wordFreq = wordFreq + int64(1)
					toInsertPage = append(toInsertPage, string(strconv.FormatInt(termId, 10)+" "+strconv.FormatInt(wordFreq, 10)))
				} else {
					toInsertPage = append(toInsertPage, string(val))
				}
			}
			if wordExists == false {
				toInsertPage = append(toInsertPage, string(strconv.FormatInt(termId, 10)+" "+strconv.FormatInt(int64(1), 10)))
			}
		} else {
			// insert a new entry
			toInsertPage = append(toInsertPage, string(strconv.FormatInt(termId, 10)+" "+strconv.FormatInt(int64(1), 10)))
		}

		err = wordDb.Update(func(tx *bolt.Tx) error {
			wordFreqBucket := tx.Bucket([]byte(wordFreqBuck))
			// fmt.Println("toInsert: ", toInsert)
			// fmt.Println("page is: ", strconv.FormatInt(pageId, 10), "word: " , word, "termId: ", termId)
			err := wordFreqBucket.Put(IntToByte(termId), StringToByte(toInsert))
			if err != nil {
				return fmt.Errorf("Error when updating wordFreqBucket: %s", err)
			}

			pageWordFreqBucket := tx.Bucket([]byte(pageWordFreqBuck))
			err = pageWordFreqBucket.Put(IntToByte(pageId), StringToByte(toInsertPage))
			if err != nil {
				return fmt.Errorf("Error when updating pageWordFreqBucket: %s", err)
			}

			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

// GetTopWords given a page id return the top 5 most frequent keywords
func GetTopWords(id int64) []string {
	exist := true
	var value []byte
	var ret []string
	var topWords []string
	err := wordDb.View(func(tx *bolt.Tx) error {
		pageWordFreqBucket := tx.Bucket([]byte(pageWordFreqBuck))
		value = pageWordFreqBucket.Get(IntToByte(id))
		if value == nil {
			exist = false
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	var a [][2]int64
	if exist == true {
		ret = ByteToString(value)
		for _, val := range ret {
			test := strings.Split(val, " ")
			keyId, err := strconv.ParseInt(test[0], 10, 64)
			if err != nil {
				fmt.Errorf("Error in wordId, converting keyId: %s", err)
			}
			freq, err := strconv.ParseInt(test[1], 10, 64)
			if err != nil {
				fmt.Errorf("Error in wordId, converting keyId: %s", err)
			}
			var tes [2]int64
			tes[0] = keyId
			tes[1] = freq
			a = append(a, tes)
		}
		sort.Slice(a, func(i, j int) bool {
			return a[i][1] > a[j][1]
		})
		// now a is sorted
		// get the top 5
		for i := 0; i < 5; i++ {
			topWords = append(topWords, string(GetWord(a[i][0])+": "+strconv.FormatInt(a[i][1], 10)))
		}
	} else {
		return nil
	}
	return topWords
}

// parses all the words given to the buckets
func parseAllWord(page *crawler.Page) {
	pageId := GetPageId(page.GetURL())
	words := page.GetKeywords()
	// first, we convert each keyword into it's wordId
	for _, word := range words {
		// check if the word is already in the database
		_ = createWordId(word)
		updateFreq(pageId, word)
	}
}

// given a url, get all the keyword statistics on that page
func GetPageKeyFreq(pageId int64) (ret []string) {
	//   pageId := GetPageId(url)
	ret = nil
	err := wordDb.View(func(tx *bolt.Tx) error {
		pageWordFreqBucket := tx.Bucket([]byte(pageWordFreqBuck))
		ret = ByteToString(pageWordFreqBucket.Get(IntToByte(pageId)))

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

func PrintWordDb() {
	err := wordDb.View(func(tx *bolt.Tx) error {
		fmt.Println("idToWordBucket")
		idToWordBucket := tx.Bucket([]byte(idToWordBuck))
		c := idToWordBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Println("key: ", ByteToInt(k), "value: ", string(v))
		}

		wordToIdBucket := tx.Bucket([]byte(wordToIdBuck))
		c = wordToIdBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Println("key: ", string(k), "value: ", ByteToInt(v))
		}

		wordFreqBucket := tx.Bucket([]byte(wordFreqBuck))
		c = wordFreqBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Println("key: ", GetWord(ByteToInt(k)), "value: ", ByteToString(v))
		}

		pageWordFreqBucket := tx.Bucket([]byte(pageWordFreqBuck))
		c = pageWordFreqBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Println("key: ", GetPageUrl(ByteToInt(k)), "value: ", ByteToString(v))
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
