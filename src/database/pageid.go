package database

import (
	"fmt"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

var pageid *bolt.DB
var firstBucket string = "pageToIdBuck"
var secondBucket string = "idToPageBuck"

// open pageid database
// initialise buckets
func OpenPageDb() {
  var err error
  pageid, err = bolt.Open("db"+string(os.PathSeparator)+"pageid.db", 0700, nil)
  if err != nil {
    log.Fatal(err)
  }

  err = pageid.Update(func(tx *bolt.Tx) error {
    pageToId, err := tx.CreateBucketIfNotExists([]byte(firstBucket))
    if err != nil {
      return fmt.Errorf("create bucket: %s", err)
    }

    // initialise the first key/value pair of firstBucket to be the total number of pages
    pageCount := pageToId.Get(intToByte(0))

    if pageCount == nil {
      pageCount = intToByte(0)
      err := pageToId.Put(intToByte(0), pageCount)
      if err != nil {
        return fmt.Errorf("Initialise pageCount error: %s", err)
      }
    }
    
    _, err = tx.CreateBucketIfNotExists([]byte(secondBucket))
    if err != nil {
      return fmt.Errorf("create bucket: %s", err)
    }

    return nil
  })
  if err != nil {
    log.Fatal(err)
  }
}

// close pageid database
func ClosePageDb() {
  pageid.Close()
}

// returns pageId of the given url, if the given page url does not exist, create a new one
func GetPageId(url string) (id int64) {
  // id of the new page, if id is -1, it means page does not exist and have to insert
  id = 0
  // count is the number of pages in the db
  var count int64 = 0

  // check if url of the page already exists first
  // if exists, change id to 
  pageid.View(func(tx *bolt.Tx) error {
    pageToId := tx.Bucket([]byte(firstBucket))
    value := pageToId.Get([]byte(url))
    // page does not exist yet
    if value == nil {
      id = -1
      count = byteToInt(pageToId.Get(intToByte(0)))
    } else {
      id = byteToInt(value)
    }

   return nil 
  })

  // if page does not exist yet, insert
  if id == -1 {
    pageid.Update(func(tx *bolt.Tx) error {
      count += 1

      // insert the new page
      pageToId := tx.Bucket([]byte(firstBucket))
      err := pageToId.Put([]byte(url), intToByte(count))
      if err != nil {
        return err
      }

      // update the count of the pages in the db
      err = pageToId.Put(intToByte(0), intToByte(count))
      if err != nil {
        return err
      }

      idToPage := tx.Bucket([]byte(secondBucket))
      err = idToPage.Put(intToByte(count), []byte(url))
      if err != nil {
        return err
      }

      return nil
    })
  }
  return
}

func GetPageUrl(id int64) (url string) {
  var value []byte
  pageid.View(func(tx *bolt.Tx) error {
    idToPage := tx.Bucket([]byte(secondBucket))
    value = idToPage.Get(intToByte(id))
    return nil
  })

  if value == nil {
    return ""
  }
  url = string(value)
  return
}
