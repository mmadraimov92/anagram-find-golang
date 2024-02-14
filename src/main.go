package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	word    = flag.String("w", "", "Input word to search for anagram")
	dict    = flag.String("d", "", "Dictionary to search from")
	charEnc = flag.String("e", "windows-1257", "Encoding of dictionary file")
)

const (
	charsNum int = 512 // Max index for char, 384 for estonian?
)

var (
	workers = 8 // Number of worker routines to spawn
)

func main() {
	start := time.Now()

	flagParse()

	var a *anagram
	a = newAnagram(dict, charEnc)
	a.findAnagram(word)
	elapsed := time.Since(start)

	fmt.Print(elapsed, ",")
	keys := reflect.ValueOf(a.result).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	fmt.Println(strings.Join(strkeys, ","))

	return
}

func flagParse() {
	flag.Parse()
	flag.Usage = func() {
		fmt.Println("Usage: anagram-find -w <word> -d <dictionary> [-e <encoding>]")
		os.Exit(2)
	}
	if *word == "" {
		flag.Usage()
	}
	if *dict == "" {
		flag.Usage()
	}
	_, ok := encodings[*charEnc]
	if !ok {
		fmt.Println("Please select correct encoding with '-e'. Check charset_table.go for reference.")
		flag.Usage()
	}
}
