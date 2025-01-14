package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {

	keyword := "fireman"

	file, err := os.Create("./hello.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content := "hello world fireman is a firefighter fireman"
	length, err := io.WriteString(file, content)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully written bytes: %d\n", length)

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return
	}
	fileSize := fileInfo.Size()

	numChunks := 4
	chunkSize := fileSize / int64(numChunks)
	overlap := int64(len(keyword) - 1)

	results := make(chan int)

	var wg sync.WaitGroup

	readFile, err := os.Open(file.Name())
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer readFile.Close()

	for i := 0; i < numChunks; i++ {
		wg.Add(1)

		startOffset := int64(i) * chunkSize
		endOffset := startOffset + chunkSize + overlap
		if endOffset > fileSize {
			endOffset = fileSize
		}

		go func(start, end int64) {
			defer wg.Done()
			findKeywordInChunk(readFile, start, end, keyword, results)
		}(startOffset, endOffset)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("Keyword found at byte offsets:")
	for result := range results {
		fmt.Println(result)
	}
}

func findKeywordInChunk(file *os.File, start, end int64, keyword string, results chan<- int) {
	buffer := make([]byte, end-start)
	_, err := file.ReadAt(buffer, start)
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading file at offset %d: %v\n", start, err)
		return
	}

	for {
		index := bytes.Index(buffer, []byte(keyword))
		if index == -1 {
			break
		}

		results <- int(start) + index

		buffer = buffer[index+len(keyword):]
		start += int64(index + len(keyword))
	}
}
