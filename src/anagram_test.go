package main

import (
	"reflect"
	"testing"
)

func BenchmarkAnagram(b *testing.B) {
	var word = "search"
	var dict = "../rockyou.txt"

	for b.Loop() {
		var a *anagram
		a = newAnagram(&dict)
		a.findAnagram(&word)
	}
}

func TestAnagram(t *testing.T) {
	tests := []struct {
		name string
		word string
		dict string
		keys []string
	}{
		{"Estonian1", "eesti", "../rockyou.txt",
			[]string{"istee", "siete", "ieste", "etsie", "tesie", "isete", "seite", "eetsi", "estie"}},
	}

	for _, tt := range tests {
		var a *anagram
		a = newAnagram(&tt.dict)
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
