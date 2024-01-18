package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// reads a file as an array of bytes
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

func decompress(data []byte) []byte {
	// adding the zlib header
	data = append(data, 0x78)
	data = append(data, 0x1)

	buffer := bytes.NewBuffer(data)

	r, err := zlib.NewReader(buffer)

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	var decompressedBytes bytes.Buffer

	io.Copy(&decompressedBytes, r)

	return decompressedBytes.Bytes()
}

func main() {
	// parsing the arguments passed to the program
	commandLineArgs := os.Args

	if len(commandLineArgs) != 2 {
		fmt.Printf("usage: %s <file>\n", filepath.Base(commandLineArgs[0]))
		os.Exit(1)
	}

	fileBytes := readFileBytes(commandLineArgs[1])

	// getting the easily accessible data
	signature := fileBytes[:3]
	compressionSignature := fileBytes[0]
	version := fileBytes[3]
	bytesSize := binary.LittleEndian.Uint32(fileBytes[4:8])

	if fileBytes[0] == 'C' {
		decompressed := decompress(fileBytes[8:])
		fmt.Printf("%s %c %d %d %x", signature, compressionSignature, version, bytesSize, decompressed[0])
	}

}
