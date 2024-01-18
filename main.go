package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Rect struct {
	Nbits uint64
	Xmin  int64
	Xmax  int64
	Ymin  int64
	Ymax  int64
}

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

func fmtBits(data []byte) []byte {
	var buf bytes.Buffer
	for _, b := range data {
		fmt.Fprintf(&buf, "%08b ", b)
	}
	buf.Truncate(buf.Len() - 1) // To remove extra space
	return buf.Bytes()
}

// gets a rect structure
func getRect(data []byte) Rect {
	var rectData Rect

	bytesSlice := data[:9]
	binaryData := string(fmtBits(bytesSlice))
	binaryData = strings.Replace(binaryData, " ", "", -1)

	rectData.Nbits, _ = strconv.ParseUint(binaryData[:5], 2, 8)
	rectData.Xmin, _ = strconv.ParseInt(binaryData[5:20], 2, 15)
	rectData.Xmax, _ = strconv.ParseInt(binaryData[20:35], 2, 15)
	rectData.Ymin, _ = strconv.ParseInt(binaryData[35:50], 2, 15)
	rectData.Ymax, _ = strconv.ParseInt(binaryData[50:65], 2, 15)

	return rectData
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

	var compressionType string

	// decompressing
	if compressionSignature == 'C' {
		compressionType = "zlib compressed"

		fileBytes = decompress(fileBytes[8:])
	} else if compressionSignature == 'F' {
		fileBytes = fileBytes[8:]
		compressionType = "uncompressed"
	} else {
		compressionType = "LZMA compressed"
	}

	getRect(fileBytes)

	fmt.Printf("Signature: %s\nCompression type: %c (%s)\nSWF version: %d\nUncompressed size (bytes): %d\n", signature, compressionSignature, compressionType, version, bytesSize)
}
