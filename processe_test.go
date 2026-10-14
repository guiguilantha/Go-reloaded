package main

import "testing"

func TestFixArticles(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "a before vowel", input: "this is a apple", want: "this is an apple"},
		{name: "a before consonant", input: "this is a banana", want: "this is a banana"},
		{name: "an before consonant", input: "this is an banana", want: "this is a banana"},
		{name: "silent h", input: "this is a hour", want: "this is an hour"},
		{name: "uppercase start", input: "A apple", want: "An apple"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fixArticles(tt.input); got != tt.want {
				t.Fatalf("fixArticles(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
