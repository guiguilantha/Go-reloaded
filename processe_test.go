package main

import "testing"

func TestFixQuote(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{name: "space before quote", in: "hello\"world\"", out: "hello \"world\""},
		{name: "space after quote", in: "\"hello\"world", out: "\"hello\" world"},
		{name: "with punctuation", in: "He said:\"hello\".", out: "He said: \"hello\"."},
		{name: "multiple quotes", in: "hello\"world\"again", out: "hello \"world\" again"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fixQuote(tt.in)
			if got != tt.out {
				t.Fatalf("fixQuote(%q) = %q, want %q", tt.in, got, tt.out)
			}
		})
	}
}
