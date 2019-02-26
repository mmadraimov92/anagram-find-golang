package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/transform"
)

var (
	word    = flag.String("w", "", "Input word to search for anagram")
	dict    = flag.String("d", "lemmad.txt", "Dictionary to search from")
	charEnc = flag.String("e", "windows-1257", "Encoding of dictionary file")
)

const (
	charsNum = 512    // Max index for char, 384 for estonian?
	workers  = 4 * 16 // Number of worker routines to spawn
)

var wordsToCompare = make(chan string, workers)
var result []string
var wg sync.WaitGroup

func main() {
	start := time.Now()
	flag.Parse()
	if *word == "" {
		fmt.Println("Please input word with '-w'")
		return
	}

	enc, ok := encodings[*charEnc]
	if !ok {
		fmt.Println("Please select correct encoding with '-e'. Check charset_table.go for reference.")
		return
	}

	startJob(&enc)
	elapsed := time.Since(start)
	fmt.Println(int64(elapsed/time.Microsecond), "µs")
	fmt.Println("Word: ", *word)
	fmt.Println("Anagrams:", result)
	return
}

func startJob(enc *encType) {
	go producer()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(enc)
	}
	wg.Wait()
}

func producer() {
	file, err := os.Open(*dict)
	defer func() {
		file.Close()
		close(wordsToCompare)
	}()
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		wordsToCompare <- line
	}
}

func worker(enc *encType) {
	defer wg.Done()

	for line := range wordsToCompare {
		wordFromDict, _, err := transform.String(*enc.e.NewDecoder(), line)
		check(err)
		if isAnagram(word, &wordFromDict) {
			result = append(result, wordFromDict)
		}
	}
}

func isAnagram(str1, str2 *string) bool {
	if len(*str1) != len(*str2) {
		return false
	}

	count := make([]int, charsNum)

	var i int
	for i = 0; i < len(*str1); {
		r1, size := utf8.DecodeRuneInString((*str1)[i:])
		count[int(unicode.ToLower(r1))]++
		i += size
	}

	for i = 0; i < len(*str2); {
		r2, size := utf8.DecodeRuneInString((*str2)[i:])
		count[int(unicode.ToLower(r2))]--
		i += size
	}

	for i = 0; i < charsNum; i++ {
		if count[i] != 0 {
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
