package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"sync"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/transform"
)

type anagram struct {
	enc            encType     // charset of dictionary file - default windows-1257
	dictionary     string      // dictionary file
	wordsToCompare chan string // channel to send dictionary words
	result         []string    // list of found anagrams
	mutex          sync.Mutex
	wg             sync.WaitGroup
}

func newAnagram(dict, charset *string) *anagram {
	var a anagram

	a.dictionary = *dict
	a.enc = encodings[*charset]
	a.wordsToCompare = make(chan string, workers)

	return &a
}

func (a *anagram) findAnagram(word *string) {
	go a.producer()

	for i := 0; i < workers; i++ {
		a.wg.Add(1)
		go a.worker(word)
	}
	a.wg.Wait()
}

func (a *anagram) producer() {
	defer func() {
		close(a.wordsToCompare)
	}()

	producers := workers

	var done = make(chan bool, producers)
	var chunks = make(chan []byte, producers)

	file, err := ioutil.ReadFile(a.dictionary)
	check(err)

	go split(file, len(file)/producers, chunks)

	for v := range chunks {
		go func(v []byte) {
			scanner := bufio.NewScanner(bytes.NewReader(v))
			for scanner.Scan() {
				a.wordsToCompare <- scanner.Text()
			}
			done <- true
		}(v)
	}

	for i := 0; i < producers; i++ {
		<-done
	}
}

func (a *anagram) worker(word *string) {
	defer a.wg.Done()

	for line := range a.wordsToCompare {
		if utf8.RuneCountInString(*word) != utf8.RuneCountInString(line) {
			continue
		}
		wordFromDict, _, err := transform.String(a.enc.e.NewDecoder(), line)
		check(err)
		if isAnagram(*word, wordFromDict) {
			a.mutex.Lock()
			a.result = append(a.result, wordFromDict)
			a.mutex.Unlock()
		}
	}
}

func isAnagram(str1, str2 string) bool {
	if len(str1) != len(str2) {
		return false
	}
	// histSize := charsNum
	histogram := make([]int, charsNum)

	for _, r1 := range str1 {
		ord := int(unicode.ToLower(r1))
		histogram[ord]++
	}

	for _, r2 := range str2 {
		ord := int(unicode.ToLower(r2))
		histogram[ord]--
	}

	for i := 0; i < charsNum; i++ {
		if histogram[i] != 0 {
			return false
		}
	}

	return true
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Split []byte array by "\n" into equal []byte arrays
func split(buf []byte, lim int, chunks chan<- []byte) {
	defer close(chunks)
	var chunk []byte
	for len(buf) >= lim {
		for i, v := range buf[lim:] {
			if v == 10 {
				chunk, buf = buf[:lim+i], buf[lim+i+1:]
				break
			}
		}
		chunks <- chunk
	}
	if len(buf) > 0 {
		chunks <- buf
	}
}
