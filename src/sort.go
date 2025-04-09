package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sort"
)

// Read a big-endian uint32 from a byte slice of length at least 4
func ReadBigEndianUint32(buffer []byte) uint32 {
	if len(buffer) < 4 {
		panic("buffer too short to read uint32")
	}
	return binary.BigEndian.Uint32(buffer[:])
}

// Write a big-endian uint32 to a byte slice of length at least 4
func WriteBigEndianUint32(buffer []byte, num uint32) {
	if len(buffer) < 4 {
		panic("buffer too short to write uint32")
	}
	binary.BigEndian.PutUint32(buffer, num)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}
	dt, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	doubleSlice := [][]byte{}
	length := len(dt)
	point := 0
	for point < length {
		curSize := ReadBigEndianUint32(dt[point : point+4])
		doubleSlice = append(doubleSlice, dt[point:point+4+int(curSize)])
		point = point + 4 + int(curSize)
	}
	sort.Slice(doubleSlice, func(i, j int) bool {
		keyi := doubleSlice[i][4:14]
		keyj := doubleSlice[j][4:14]
		for i := 0; i < 10; i++ {
			if keyi[i] != keyj[i] {
				return keyi[i] < keyj[i]
			}
		}
		return true
	})
	var buffer bytes.Buffer

	for _, each := range doubleSlice {
		buffer.Write(each)
	}
	err = os.WriteFile(os.Args[2], buffer.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
	}

}
