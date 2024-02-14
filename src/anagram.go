package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
	"unicode"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

type anagram struct {
	enc        encType             // charset of dictionary file - default windows-1257
	dictionary string              // dictionary file
	result     map[string]struct{} // list of found anagrams
	mutex      sync.Mutex
	wg         sync.WaitGroup
	dec        *encoding.Decoder
}

func newAnagram(dict, charset *string) *anagram {
	var a anagram

	a.dictionary = *dict
	a.enc = encodings[*charset]
	a.result = make(map[string]struct{})
	dec := a.enc.e.NewDecoder()
	a.dec = dec

	return &a
}

func (a *anagram) findAnagram(word *string) {
	a.wg.Add(1)
	go a.worker(word)
	a.wg.Wait()
}

func (a *anagram) worker(word *string) {
	defer a.wg.Done()

	var done = make(chan bool, workers)
	var chunks = make(chan []byte, workers)

	content, err := os.ReadFile(a.dictionary)
	checkErr(err)

	go split(content, len(content)/workers, chunks)

	for v := range chunks {
		go func(v []byte) {
			defer func() {
				done <- true
			}()
			reader := bufio.NewReader(bytes.NewReader(v))
			for {
				line, _, err := reader.ReadLine()
				if err == io.EOF {
					break
				}
				if len([]byte(*word)) != len(line) {
					continue
				}
				wordFromDict, _, err := transform.Bytes(a.dec, line)
				checkErr(err)
				if isAnagram(*word, string(wordFromDict)) {
					a.mutex.Lock()
					a.result[string(wordFromDict)] = struct{}{}
					a.mutex.Unlock()
				}
			}
		}(v)
	}

	for i := 0; i < workers; i++ {
		<-done
	}

}

func isAnagram(word, fromDict string) bool {
	histogram := make([]int8, charsNum)

	for _, r1 := range word {
		ord := int(unicode.ToLower(r1))
		histogram[ord]++
	}

	for _, r2 := range fromDict {
		ord := int(unicode.ToLower(r2))
		if ord > charsNum {
			return false
		}
		histogram[ord]--
	}

	for i := 0; i < charsNum; i++ {
		if histogram[i] != 0 {
			return false
		}
	}

	return true
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

const newlineASCII = 10

// Split []byte array by "\n" into equal []byte arrays
func split(data []byte, bytesPerWorker int, chunks chan<- []byte) {
	defer close(chunks)
	var chunk []byte
	for len(data) > bytesPerWorker {
		for i, v := range data[bytesPerWorker:] {
			if v == newlineASCII {
				chunk, data = data[:bytesPerWorker+i], data[bytesPerWorker+i+1:]
				break
			}
		}
		chunks <- chunk
	}
	if len(data) > 0 {
		chunks <- data
	}
}
