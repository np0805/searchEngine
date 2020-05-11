package database

import (
	"fmt"
	"log"
	"os"
  "../crawler"

	bolt "go.etcd.io/bbolt"
)

var pageId *bolt.DB
var pageToIdBuck string = "pageToIdBuck"
var idToPageBuck string = "idToPageBuck"

// open pageId database
// initialise buckets
func openPageIdDb() {
  var err error
  pageId, err = bolt.Open("db" + string(os.PathSeparator) + "pageId.db", 0700, nil)
  if err != nil {
    log.Fatal(err)
  }

  err = pageId.Update(func(tx *bolt.Tx) error {
    pageToId, err := tx.CreateBucketIfNotExists([]byte(pageToIdBuck))
    if err != nil {
      return fmt.Errorf("pageId create bucket: %s", err)
    }

    // initialise the first key/value pair of firstBucket to be the total number of pages
    pageCount := pageToId.Get(IntToByte(0))

    if pageCount == nil {
      pageCount = IntToByte(0)
      err := pageToId.Put(IntToByte(0), pageCount)
      if err != nil {
        return fmt.Errorf("pageId initialise pageCount error: %s", err)
      }
    }
    
    _, err = tx.CreateBucketIfNotExists([]byte(idToPageBuck))
    if err != nil {
      return fmt.Errorf("pageId create bucket: %s", err)
    }

    return nil
  })
  if err != nil {
    log.Fatal(err)
  }
}

// close pageId database
func closePageIdDb() {
  pageId.Close()
}

// returns pageId of the given url, if the given page url does not exist, create a new one
func GetPageId(url string) (id int64) {
  // id of the new page, if id is -1, it means page does not exist and have to insert
  id = 0
  // count is the number of pages in the db
  var count int64 = 0

  // check if url of the page already exists first
  // if exists, change id to 
  pageId.View(func(tx *bolt.Tx) error {
    pageToId := tx.Bucket([]byte(pageToIdBuck))
    value := pageToId.Get([]byte(url))
    // page does not exist yet
    if value == nil {
      id = -1
      count = ByteToInt(pageToId.Get(IntToByte(0)))
    } else {
      id = ByteToInt(value)
    }

   return nil 
  })

  // if page does not exist yet, insert
  if id == -1 {
    pageId.Update(func(tx *bolt.Tx) error {
      count += 1

      // insert the new page
      pageToId := tx.Bucket([]byte(pageToIdBuck))
      err := pageToId.Put([]byte(url), IntToByte(count))
      if err != nil {
        return err
      }

      // update the count of the pages in the db
      err = pageToId.Put(IntToByte(0), IntToByte(count))
      if err != nil {
        return err
      }

      idToPage := tx.Bucket([]byte(idToPageBuck))
      err = idToPage.Put(IntToByte(count), []byte(url))
      if err != nil {
        return err
      }

      return nil
    })
  }
  return
}

// return number of pages in the pageIdDb, if empty, return 0
func GetPageNumber() (ret int64) {
  ret = int64(0)
  err := pageId.View(func(tx *bolt.Tx) error {
    pageToId := tx.Bucket([]byte(pageToIdBuck))
    value := pageToId.Get(IntToByte(int64(0)))
    if value != nil {
      ret = ByteToInt(value)
    }

    return nil
  })
  if err != nil {
    log.Fatal(err)
  }
  return ret
}

func GetPageUrl(id int64) (url string) {
  var value []byte
  pageId.View(func(tx *bolt.Tx) error {
    idToPage := tx.Bucket([]byte(idToPageBuck))
    value = idToPage.Get(IntToByte(id))
    return nil
  })

  if value == nil {
    return ""
  }
  url = string(value)
  return
}

// given a map of pages, parse all the child pages in each parent pages to get their pageId
func parseAllChild(parent *crawler.Page) {
  // fmt.Println(parent.GetChildrenURL())
  for _, child := range parent.GetChildrenURL() {
    _ = GetPageId(child)
  }
}

func PrintPageIdDb() {
  pageId.View(func(tx *bolt.Tx) error {
    fmt.Println("PAGE_TO_ID BUCKET")
    pageToId := tx.Bucket([]byte(pageToIdBuck))
    c := pageToId.Cursor()

    k, v := c.First()
    key := ByteToInt(k)
    value := ByteToInt(v)
    fmt.Println("key: ", key, "value: ", value)
    k, v = c.Next()
    for ; k != nil; k, v = c.Next() {
      key := string(k)
      value := ByteToInt(v)
      fmt.Println("key: ", key, "value: ", value)
    }

    fmt.Println("ID_TO_PAGE BUCKET")
    idToPage := tx.Bucket([]byte(idToPageBuck))
    c = idToPage.Cursor()
    for k, v = c.First(); k != nil; k, v = c.Next() {
      key := ByteToInt(k)
      value := string(v)
      fmt.Println("key: ", key, "value: ", value)
    }


    return nil
  })
}
