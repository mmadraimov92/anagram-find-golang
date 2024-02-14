package main

import (
	"reflect"
	"testing"
)

func BenchmarkAnagram(b *testing.B) {
	var word = "eesti"
	var dict = "../rockyou.txt"
	var charEnc = "windows-1257"

	for n := 0; n < b.N; n++ {
		var a *anagram
		a = newAnagram(&dict, &charEnc)
		a.findAnagram(&word)
	}
}

func TestAnagram(t *testing.T) {
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
