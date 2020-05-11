package database

import (
	"encoding/binary"
	"fmt"
	"math"

	"../crawler"
)

const maxInt32 = 1<<(32-1) - 1

func OpenAllDb() {
	openPageIdDb()
	openPageInfoDb()
	openWordDb()
}

func CloseAllDb() {
	closePageIdDb()
	closePageInfoDb()
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func IntToByte(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func ByteToInt(b []byte) (id int64) {
	id = int64(binary.LittleEndian.Uint64(b))
	return id
}

func writeLen(b []byte, l int) []byte {
	if 0 > l || l > maxInt32 {
		panic("writeLen: invalid length")
	}
	var lb [4]byte
	binary.BigEndian.PutUint32(lb[:], uint32(l))
	return append(b, lb[:]...)
}

func readLen(b []byte) ([]byte, int) {
	if len(b) < 4 {
		panic("readLen: invalid length")
	}
	l := binary.BigEndian.Uint32(b)
	if l > maxInt32 {
		panic("readLen: invalid length")
	}
	return b[4:], int(l)
}
func ByteToString(b []byte) []string {
	b, ls := readLen(b)
	s := make([]string, ls)
	for i := range s {
		b, ls = readLen(b)
		s[i] = string(b[:ls])
		b = b[ls:]
	}
	return s
}

func StringToByte(s []string) []byte {
	var b []byte
	b = writeLen(b, len(s))
	for _, ss := range s {
		b = writeLen(b, len(ss))
		b = append(b, ss...)
	}
	return b
}

// given a map of pages, parse all the parent pages to get their pageId
func ParseAllPages(pages *map[string]*crawler.Page) {
	for _, page := range *pages {
		fmt.Println("page: ", page.GetTitle())
		// fmt.Println("keywords: ", page.GetKeywords())
		_ = GetPageId(page.GetURL())
		parseAllChild(page)
		parseAllInfo(page)
		// fmt.Println("keywords: ", page.GetKeywords())
		// PrintPageInfoDb()
		parseAllWord(page)
		// PrintWordDb()
		// fmt.Println(GetPageKeyFreq(page.GetURL()))
	}
}
