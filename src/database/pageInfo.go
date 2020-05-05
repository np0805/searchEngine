package database

import (
	"fmt"
	"log"
	"os"

	"../crawler"

	bolt "go.etcd.io/bbolt"
)

var pageInfo *bolt.DB
var pageInfoBuck string = "pageInfoBuck"
var parentChildBuck string = "parentChildBuck"
var childParentBuck string = "childParentBuck"

func openPageInfoDb() {
  var err error
  pageInfo, err = bolt.Open("db"+string(os.PathSeparator)+"pageInfo.db", 0700, nil)

  if err != nil {
    log.Fatal(err)
  }

  err = pageInfo.Update(func(tx *bolt.Tx) error {
    _, err := tx.CreateBucketIfNotExists([]byte(pageInfoBuck))
    if err != nil {
      return fmt.Errorf("pageInfo create first bucket: %s", err)
    }

    _, err = tx.CreateBucketIfNotExists([]byte(parentChildBuck))
    if err != nil {
      return fmt.Errorf("pageInfo create second bucket: %s", err)
    }

    _, err = tx.CreateBucketIfNotExists([]byte(childParentBuck))
    if err != nil {
      return fmt.Errorf("pageInfo create third bucket: %s", err)
    }

    return nil
  })

  if err != nil {
    log.Fatal(err)
  }
}

func closePageInfoDb() {
  pageInfo.Close()
}

// func getAllChild(url string) (allChild []string) {
// 
// }

// parse all the info on the given page, including all it's child links
func parseAllInfo(page *crawler.Page) {
  var info []string
  info = append(info, page.GetTitle())
  info = append(info, page.GetURL())
  info = append(info, page.GetLastModified())
  info = append(info, page.GetSize())

  // value := &bytes.Buffer{}
  // enc := gob.NewEncoder(value)
  // enc.Encode(info)

  value:= StringToByte(info)
  err := pageInfo.Update(func(tx *bolt.Tx) error {
    // store the info in pageInfoBucket
    pageInfoBucket := tx.Bucket([]byte(pageInfoBuck))
    pageId := IntToByte(GetPageId(page.GetURL()))
    err := pageInfoBucket.Put(pageId, value)
    if err != nil {
      return fmt.Errorf("Error in pageInfo: parseAllInfo error: %s", err)
    }
    
    return nil
  })
  if err != nil {
    log.Fatal(err)
  }
}

// print pageInfoDb in human readable format
func PrintPageInfoDb() {
  pageInfo.View(func(tx *bolt.Tx) error {
    fmt.Println("PAGE_INFO BUCKET")
    pageInfoBucket := tx.Bucket([]byte(pageInfoBuck))
    c := pageInfoBucket.Cursor()

    for k,v := c.First(); k != nil; k, v = c.Next() {
      key := ByteToInt(k)
      //var value string = string(v)
      value := ByteToString(v)
      fmt.Println("key: ", key, "value: ", value)
      // fmt.Println("key: ", key)
    }

    return nil
  })
}
