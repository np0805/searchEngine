package database

import "encoding/binary"

func OpenAllDb() {
  openPageDb()
}

func intToByte(i int64) []byte {
  b := make([]byte, 8)
  binary.LittleEndian.PutUint64(b, uint64(i))
  return b
}

func byteToInt(b []byte) (id int64) {
  id = int64(binary.LittleEndian.Uint64(b))
  return id
}
