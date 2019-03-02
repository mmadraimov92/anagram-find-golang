package main

import (
	"testing"
)

func BenchmarkAnagram(b *testing.B) {
	var word = "hello"
	var dict = "../lemmad.txt"
	var charEnc = "windows-1257"
	b.N = 30

	for n := 0; n < b.N; n++ {
		var a *anagram
		a = newAnagram(&dict, &charEnc)
		a.findAnagram(&word)
	}
}
