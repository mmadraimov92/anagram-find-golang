package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"launchpad.net/gommap"
)

type anagram struct {
	enc            encType             // charset of dictionary file - default windows-1257
	dictionary     string              // dictionary file
	wordsToCompare chan []byte         // channel to send dictionary words
	result         map[string]struct{} // list of found anagrams
	mutex          sync.Mutex
	wg             sync.WaitGroup
	dec            *encoding.Decoder
}

func newAnagram(dict, charset *string) *anagram {
	var a anagram

	a.dictionary = *dict
	a.enc = encodings[*charset]
	a.wordsToCompare = make(chan []byte, workers*4)
	a.result = make(map[string]struct{})
	dec := a.enc.e.NewDecoder()
	a.dec = dec

	return &a
}

func (a *anagram) findAnagram(word *string) {
	go a.producer(word)

	for i := 0; i < workers; i++ {
		a.wg.Add(1)
		go a.worker(word)
	}
	a.wg.Wait()
}

func (a *anagram) producer(word *string) {
	defer func() {
		close(a.wordsToCompare)
	}()

	producers := workers

	var done = make(chan bool, producers)
	var chunks = make(chan []byte, producers)

	file, err := os.Open(a.dictionary)
	check(err)
	defer file.Close()
	mmap, err := gommap.Map(file.Fd(), gommap.PROT_READ, gommap.MAP_PRIVATE)
	check(err)

	go split(mmap, len(mmap)/producers, chunks)

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

				if utf8.RuneCountInString(*word) != len(line) {
					continue
				}
				a.wordsToCompare <- line
			}
		}(v)
	}

	for i := 0; i < producers; i++ {
		<-done
	}
}

func (a *anagram) worker(word *string) {
	defer a.wg.Done()

	for line := range a.wordsToCompare {
		wordFromDict, _, err := transform.Bytes(a.dec, line)
		check(err)
		if isAnagram(*word, string(wordFromDict)) {
			a.mutex.Lock()
			a.result[string(wordFromDict)] = struct{}{}
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
	for len(buf) > lim {
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
