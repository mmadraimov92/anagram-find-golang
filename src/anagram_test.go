package main

import (
	"reflect"
	"testing"
)

func BenchmarkAnagram(b *testing.B) {
	var word = "eesti"
	var dict = "../lemmad.txt"
	var charEnc = "windows-1257"

	for n := 0; n < b.N; n++ {
		var a *anagram
		a = newAnagram(&dict, &charEnc)
		a.findAnagram(&word)
	}
}

func TestIntegration(t *testing.T) {
	tests := []struct {
		name    string
		word    string
		dict    string
		charEnc string
		keys    []string
	}{
		{"Estonian1", "dais", "../lemmad.txt", "windows-1257", []string{"AIDS"}},
		{"Estonian2", "eesti", "../lemmad.txt", "windows-1257", []string{"eetsi", "eesti", "eiste"}},
		{"Estonian3", "tžoržet", "../lemmad.txt", "windows-1257", []string{"žoržett"}},
	}

	for _, tt := range tests {
		var a *anagram
		a = newAnagram(&tt.dict, &tt.charEnc)
		a.findAnagram(&tt.word)
		func() {
			keys := reflect.ValueOf(a.result).MapKeys()
			strkeys := make([]string, len(keys))
			for i := 0; i < len(keys); i++ {
				strkeys[i] = keys[i].String()
			}
			if len(strkeys) != len(tt.keys) {
				t.Error("Wrong number of items")
				t.Error("Expected:", tt.keys, "got:", strkeys)
			}
			for _, key := range tt.keys {
				if _, ok := a.result[key]; !ok {
					t.Error("Missing: ", key)
				}
			}
		}()
	}
}

func TestIsAnagram(t *testing.T) {
	tests := []struct {
		name   string
		word1  string
		word2  string
		answer bool
	}{
		{"Same words", "hello", "hello", true},
		{"Different length words", "hell", "ohell", false},
		{"One letter", "a", "a", true},
		{"One duplicate letter", "a", "aa", false},
		{"Spaces", "a ab cd", "  aabcd", true},
		{"Upper case", "cAmEl", "lemaC", true},
		{"Signs", "v$rt!on-suv!", "!!v$rtonsuv-", true},
		{"Estonian1", "kilööab", "öailökb", true},
		{"Estonian2", "ŠŽÕÄÖÜ", "õšžäöü", true},
		{"Empty string", "a", "", false},
		{"Numbers", "12345", "54321", true},
		// {"Chinese", "喜欢", "欢喜", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // run sub-tests in parallel
			result := isAnagram(tt.word1, tt.word2)
			if tt.answer != result {
				t.Error("Expected:", tt.answer, "got:", result)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name   string
		buf    []byte
		lim    int
		answer [][]byte
	}{
		{"Divide in half", []byte("word1\nword2\n"), 5, [][]byte{[]byte("word1"), []byte("word2")}},
		{"Small chunk size", []byte("word1\nword2\n"), 2, [][]byte{[]byte("word1"), []byte("word2")}},
		{"Long word", []byte("word1\nlongword2\n"), 7, [][]byte{[]byte("word1\nlongword2")}},
		{"Limit on new line", []byte("word1\nlongword2\n"), 6, [][]byte{[]byte("word1\nlongword2")}},
		{"Sample", []byte("word1\nword2\nword3\n"), 6, [][]byte{[]byte("word1\nword2"), []byte("word3\n")}},
		{"One word, exact limit", []byte("word1\n"), 5, [][]byte{[]byte("word1")}},
		{"One word, big limit", []byte("word1\n"), 10, [][]byte{[]byte("word1\n")}},
		{"One word, small limit", []byte("word1\n"), 1, [][]byte{[]byte("word1")}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // run sub-tests in parallel
			result := func() [][]byte {
				chunks := make(chan []byte)
				var got [][]byte
				go split(tt.buf, tt.lim, chunks)
				for chunk := range chunks {
					got = append(got, chunk)
				}
				return got
			}()

			if !reflect.DeepEqual(result, tt.answer) {
				t.Error("Expected:", tt.answer, "got:", result)
			}
		})
	}
}
