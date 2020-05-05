package database

import "encoding/binary"

func OpenAllDb() {
  openPageDb()
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
