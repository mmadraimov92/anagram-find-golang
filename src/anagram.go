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

	file, err := os.Open(a.dictionary)
	check(err)
	defer file.Close()
	mmap, err := gommap.Map(file.Fd(), gommap.PROT_READ, gommap.MAP_PRIVATE)
	check(err)

	go split(mmap, len(mmap)/workers, chunks)

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
				wordFromDict, _, err := transform.Bytes(a.dec, line)
				check(err)
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

func isAnagram(str1, str2 string) bool {
	if len(str1) != len(str2) {
		return false
	}
	sourceMap := make(map[rune]int)
	compareMap := make(map[rune]int)

	for _, r1 := range str1 {
		_, ok := sourceMap[unicode.ToLower(r1)]
		if !ok {
			sourceMap[unicode.ToLower(r1)] = 0
		} else {
			sourceMap[unicode.ToLower(r1)]++
		}
	}

	for _, r2 := range str2 {
		_, ok := compareMap[unicode.ToLower(r2)]
		if !ok {
			compareMap[unicode.ToLower(r2)] = 0
		} else {
			compareMap[unicode.ToLower(r2)]++
		}
	}

	for k, v := range sourceMap {
		if w, ok := compareMap[k]; !ok || v != w {
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
