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

func printHeader(fileBytes []byte) {
	// getting the easily accessible data
	signature := fileBytes[:3]
	compressionSignature := fileBytes[0]
	version := fileBytes[3]
	bytesSize := binary.LittleEndian.Uint32(fileBytes[4:8])

	var compressionType string

	// decompressing
	if compressionSignature == 'C' {
		compressionType = "zlib compressed"

		decompressed := decompress(fileBytes[8:])
		fileBytes = append(fileBytes[:8], decompressed...)
	} else if compressionSignature == 'F' {
		compressionType = "uncompressed"
	}

	// ensuring the signature is valid
	if string(signature) != "FWS" && string(signature) != "CWS" {
		fmt.Printf("Unknown signature: %s. Expected 'FWS' or 'CWS'.\n", string(signature))
		os.Exit(1)
	}

	frameSize := getRect(fileBytes[8:])
	frameRate := fileBytes[18]
	frameCount := fileBytes[19]

	fmt.Printf("Signature: %s\nCompression type: %c (%s)\nSWF version: %d\nUncompressed size (bytes): %d\nFrame size (twips):\n  N-bits: %d\n  X minimum: %d\n  X maximum: %d\n  Y minimum: %d\n  Y maximum: %d\nFramerate: %dfps\nFrame count: %d\n", signature, compressionSignature, compressionType, version, bytesSize, frameSize.Nbits, frameSize.Xmin, frameSize.Xmax, frameSize.Ymin, frameSize.Ymax, frameRate, frameCount)
}

func decompressFile(fileBytes []byte, fileName string) {
	if fileBytes[0] != 'C' {
		fmt.Printf("file already decompressed!")
		return
	}

	pwd, _ := os.Getwd()
	outputFilename := pwd + strings.Trim(fileName, ".swf") + "_decompiled.swf"
	file, err := os.Create(outputFilename)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	decompressed := decompress(fileBytes[8:])
	fileBytes = append(fileBytes[:8], decompressed...)
	fileBytes[0] = 'F'

	file.Write(fileBytes)
}

func main() {
	// parsing the arguments passed to the program
	commandLineArgs := os.Args

	if len(commandLineArgs) < 2 {
		fmt.Printf("usage: %s -d <file>\n", filepath.Base(commandLineArgs[0]))
		os.Exit(1)
	}

	if len(commandLineArgs) == 2 {
		fileBytes := readFileBytes(commandLineArgs[1])
		printHeader(fileBytes)
	} else if commandLineArgs[1] == "-d" {
		fileBytes := readFileBytes(commandLineArgs[2])
		decompressFile(fileBytes, commandLineArgs[2])
	} else {
		fmt.Printf("unknown argument: %s\n", commandLineArgs[1])
		os.Exit(1)
	}
}
