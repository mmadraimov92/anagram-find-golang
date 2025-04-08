package main

import (
	"os"
	"runtime"
	"sync"

	"golang.org/x/exp/mmap"
)

const (
	returnASCII  byte = '\r'
	newlineASCII byte = '\n'

	maxByte int = 256
)

var workers = runtime.NumCPU() // Number of worker routines to spawn
// var workers = 32 // Number of worker routines to spawn

type anagram struct {
	dictionary    string              // dictionary file
	result        map[string]struct{} // list of found anagrams
	mutex         sync.Mutex
	wg            sync.WaitGroup
	word          []byte
	wordLen       int
	wordHistogram []uint8
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
	a.wordHistogram = make([]uint8, maxByte)
	for _, r1 := range a.word {
		ord := uint8(r1)
		a.wordHistogram[ord]++
	}

	a.start()
}

func (a *anagram) start() {
	file, err := os.Open(a.dictionary)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	reader, err := mmap.Open(a.dictionary)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	data := make([]byte, stat.Size())
	n, err := reader.ReadAt(data, 0)
	if err != nil {
		panic(err)
	}

	a.split(data[:n], len(data[:n])/workers)
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

	histogram := make([]uint8, maxByte)
	wordBuffer := make([]byte, 50)

	var offset int
	for i, b := range chunk {
		if b == newlineASCII {
			wordLen := i - offset
			if wordLen == a.wordLen {
				copy(wordBuffer[:wordLen], chunk[offset:i])
				copy(histogram, a.wordHistogram)

				isAnagram := true
				for j := offset; j < i; j++ {
					ord := uint8(chunk[j])
					histogram[ord]--
					if histogram[ord] > 127 {
						isAnagram = false
						break
					}
				}

				if isAnagram && isHistogramEmpty(histogram) {
					a.mutex.Lock()
					a.result[string(wordBuffer[:wordLen])] = struct{}{}
					a.mutex.Unlock()
				}
			}
			offset = i + 1
		}
	}
}

func isHistogramEmpty(histogram []uint8) bool {
	for i := range maxByte {
		if histogram[i] != 0 {
			return false
		}
	}
	return true
}
