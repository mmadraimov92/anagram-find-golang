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

func TestIsAnagram(t *testing.T) {
	tests := []struct {
		msg    string
		word1  string
		word2  string
		answer bool
	}{
		{"Same words", "hello", "ehllo", true},
		{"Different length words", "hell", "ohell", false},
		{"One letter", "a", "a", true},
		{"One duplicate letter", "a", "aa", false},
		{"Spaces", "a ab cd", "  aabcd", true},
		{"Upper case", "cAmEl", "lemaC", true},
		{"Signs", "v$rt!on-suv!", "!!v$rtonsuv-", true},
		{"Estonian", "kilööab", "öailökb", true},
		// {"Chinese", "喜欢", "欢喜", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			t.Parallel() // run sub-tests in parallel
			result := isAnagram(tt.word1, tt.word2)
			if tt.answer != result {
				t.Error("Expected:", tt.answer, "got:", result)
			}
		})
	}
}
