package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"

	"golang.org/x/text/transform"
)

var (
	word    = flag.String("w", "", "Input word to search for anagram")
	dict    = flag.String("d", "", "Dictionary to search from")
	charEnc = flag.String("e", "windows-1257", "Encoding of dictionary file")
)

const (
	charsNum = 512    // Max index for char, 384 for estonian?
	workers  = 4 * 16 // Number of worker routines to spawn
)

type anagram struct {
	mainWord       string      // Main word to find anagrams of
	enc            encType     // charset of dictionary file - default windows-1257
	dictionary     string      // dictionary file
	wordsToCompare chan string // channel to send dictionary words
	result         []string    // list of found anagrams
	mutex          sync.Mutex
	wg             sync.WaitGroup
}

func main() {
	start := time.Now()

	var anagram anagram
	anagram.initializeAnagram()

	anagram.findAnagram()
	elapsed := time.Since(start)
	fmt.Print(int64(elapsed/time.Microsecond), ",")
	fmt.Println(strings.Join(anagram.result, ","))
	return
}

func (a *anagram) initializeAnagram() {
	flag.Usage = func() {
		fmt.Println("Usage: anagram-find -w <word> -d <dictionary> [-e <encoding>]")
		os.Exit(2)
	}
	flag.Parse()
	if *word == "" {
		flag.Usage()
	}
	a.mainWord = *word

	if *dict == "" {
		flag.Usage()
	}
	a.dictionary = *dict

	var ok bool
	a.enc, ok = encodings[*charEnc]
	if !ok {
		fmt.Println("Please select correct encoding with '-e'. Check charset_table.go for reference.")
		flag.Usage()
	}

	a.wordsToCompare = make(chan string, workers)
}

func (a *anagram) findAnagram() {
	go a.producer()

	for i := 0; i < workers; i++ {
		a.wg.Add(1)
		go a.worker()
	}
	a.wg.Wait()
}

func (a *anagram) producer() {
	file, err := os.Open(a.dictionary)
	check(err)
	defer func() {
		file.Close()
		close(a.wordsToCompare)
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		a.wordsToCompare <- scanner.Text()
	}

}

func (a *anagram) worker() {
	defer a.wg.Done()

	for line := range a.wordsToCompare {
		wordFromDict, _, err := transform.String(a.enc.e.NewDecoder(), line)
		check(err)
		if isAnagram(a.mainWord, wordFromDict) {
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

	count := make([]int, charsNum)

	for _, r1 := range str1 {
		count[int(unicode.ToLower(r1))]++
	}

	for _, r2 := range str2 {
		count[int(unicode.ToLower(r2))]--
	}

	for i := 0; i < charsNum; i++ {
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
