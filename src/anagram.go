package main

import (
	"os"
	"sync"
)

const (
	charsNum          = 384 // Max index for char, 384 for estonian?
	workers           = 8   // Number of worker routines to spawn
	returnASCII  byte = '\r'
	newlineASCII byte = '\n'
)

type anagram struct {
	dictionary string              // dictionary file
	result     map[string]struct{} // list of found anagrams
	mutex      sync.Mutex
	wg         sync.WaitGroup
	word       []byte
	wordLen    int
}

func newAnagram(dict *string) *anagram {
	var a anagram

	a.dictionary = *dict
	a.result = make(map[string]struct{})

	return &a
}

func (a *anagram) findAnagram(word *string) {
	a.word = []byte(*word)
	a.wordLen = len(*word)
	a.start()
}

func (a *anagram) start() {
	content, err := os.ReadFile(a.dictionary)
	checkErr(err)

	a.split(content, len(content)/workers)
}

func isAnagram(word, fromDict []byte) bool {
	histogram := make([]int8, charsNum)

	for _, r1 := range word {
		ord := int(r1)
		histogram[ord]++
	}

	for _, r2 := range fromDict {
		ord := int(r2)
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

// Split []byte array by "\n" into equal []byte arrays
func (a *anagram) split(data []byte, bytesPerWorker int) {
	var chunk []byte
	for len(data) > bytesPerWorker {
		for i, v := range data[bytesPerWorker:] {
			if v == newlineASCII {
				chunk, data = data[:bytesPerWorker+i], data[bytesPerWorker+i+1:]
				break
			}
		}
		a.wg.Add(1)
		go a.process(chunk)
	}
	if len(data) > 0 {
		a.wg.Add(1)
		go a.process(data)
	}
	a.wg.Wait()
}

func (a *anagram) process(chunk []byte) {
	defer a.wg.Done()

	var line []byte
	var offset int
	for i, b := range chunk {
		if b == newlineASCII {
			line = chunk[offset:i]
			if chunk[i-1] == returnASCII {
				line = line[:len(line)-1]
			}
			offset = i + 1
			if a.wordLen != len(line) {
				continue
			}
			if isAnagram(a.word, line) {
				a.mutex.Lock()
				a.result[string(line)] = struct{}{}
				a.mutex.Unlock()
			}
		}
	}
}
