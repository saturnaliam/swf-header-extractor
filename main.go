package main

import (
	"encoding/binary"
	"io"
	"os"
)

func readFileBytes(filename string) []byte {
	file, err := os.Open(filename)
	var fileBytes []byte

	if err != nil {
		panic(err)
	}

	defer file.Close()

	for {
		var byteRead byte

		err = binary.Read(file, binary.LittleEndian, &byteRead)

		if err != nil {
			if err == io.EOF {
				break
			}
		}

		fileBytes = append(fileBytes, byteRead)
	}

	return fileBytes
}

func main() {
	readFileBytes("pvz_9_15.swf")
}
