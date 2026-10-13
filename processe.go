package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	markerRegex = regexp.MustCompile(`\((hex|bin|up|low|cap)(?:,\s*(\d+))?\)`)
	punctBefore = regexp.MustCompile(`\s+([.,!?:;]+)`)
	punctAfter  = regexp.MustCompile(`([.,!?:;]+)([^\s.,!?:;])`)
	quoteRegex  = regexp.MustCompile(`"([^"]*)"`)
)

func process(content string) string {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		line = processMarkers(line)
		line = fixQuote(line)
		line = fixPunct(line)
		lines[i] = line
	}

	return strings.Join(lines, "\n")
}

func processMarkers(line string) string {
	line = markerRegex.ReplaceAllStringFunc(line, func(match string) string {
		return " " + strings.ReplaceAll(match, " ", "") + " "
	})

	words := strings.Fields(line)
	out := []string{}

	for _, word := range words {
		match := markerRegex.FindStringSubmatch(word)

		if match == nil {
			out = append(out, word)
			continue
		}

		cmd := match[1]
		n := 1

		if match[2] != "" {
			n, _ = strconv.Atoi(match[2])
		}

		for i := len(out) - 1; i >= 0 && n > 0; i-- {
			if !isWord(out[i]) {
				continue
			}

			out[i] = transform(cmd, out[i])
			n--
		}
	}

	return strings.Join(out, " ")
}

func isWord(s string) bool {
	return regexp.MustCompile(`[A-Za-zÀ-ÿ]`).MatchString(s)
}

func transform(cmd, word string) string {
	switch cmd {
	case "up":
		return strings.ToUpper(word)

	case "low":
		return strings.ToLower(word)

	case "cap":
		runes := []rune(strings.ToLower(word))

		if len(runes) > 0 {
			runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
		}

		return string(runes)

	case "hex":
		n, err := strconv.ParseInt(word, 16, 64)
		if err != nil {
			return word
		}

		return strconv.FormatInt(n, 10)

	case "bin":
		n, err := strconv.ParseInt(word, 2, 64)
		if err != nil {
			return word
		}

		return strconv.FormatInt(n, 10)
	}

	return word
}

func fixQuote(line string) string {
	runes := []rune(line)
	var out strings.Builder
	inQuote := false
	prev := rune(0)

	for i, r := range runes {
		if r != '"' {
			out.WriteRune(r)
			prev = r
			continue
		}

		if !inQuote {
			if prev != 0 && !unicode.IsSpace(prev) {
				out.WriteRune(' ')
			}
			out.WriteRune(r)
			inQuote = true
			prev = r
			continue
		}

		out.WriteRune(r)
		inQuote = false
		if i+1 < len(runes) {
			next := runes[i+1]
			if !unicode.IsSpace(next) && !isQuoteSeparator(next) {
				out.WriteRune(' ')
			}
		}
		prev = r
	}

	return out.String()
}

func isQuoteSeparator(r rune) bool {
	switch r {
	case '.', ',', '!', '?', ':', ';':
		return true
	default:
		return false
	}
}

func fixPunct(line string) string {
	line = punctBefore.ReplaceAllString(line, "$1")
	line = punctAfter.ReplaceAllString(line, "$1 $2")

	return line
}
