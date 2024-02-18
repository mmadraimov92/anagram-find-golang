package main

import (
	"os"
	"runtime"
	"sync"
)

const (
	returnASCII  byte = '\r'
	newlineASCII byte = '\n'
)

var workers = runtime.NumCPU() // Number of worker routines to spawn

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
	if err != nil {
		panic(err)
	}

	a.split(content, len(content)/workers)
}

func isAnagram(word, fromDict []byte) bool {
	histogram := make([]int8, 256)

	for _, r1 := range word {
		ord := uint8(r1)
		histogram[ord]++
	}

	for _, r2 := range fromDict {
		ord := uint8(r2)
		histogram[ord]--
	}

	for i := 0; i < 256; i++ {
		if histogram[i] != 0 {
			return false
		}
	}

	return true
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

	var offset int
	for i, b := range chunk {
		if b == newlineASCII {
			if a.wordLen == len(chunk[offset:i]) && isAnagram(a.word, chunk[offset:i]) {
				a.mutex.Lock()
				a.result[string(chunk[offset:i])] = struct{}{}
				a.mutex.Unlock()
			}
			offset = i + 1
		}
	}
}
