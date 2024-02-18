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
	word = flag.String("w", "", "Input word to search for anagram")
	dict = flag.String("d", "", "Dictionary to search from")
)

func main() {
	start := time.Now()

	flagParse()

	var a *anagram
	a = newAnagram(dict)
	a.findAnagram(word)
	elapsed := time.Since(start)

	fmt.Print(elapsed, ",")
	keys := reflect.ValueOf(a.result).MapKeys()
	result := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		result[i] = keys[i].String()
	}
	fmt.Println(strings.Join(result, ","))

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
}
