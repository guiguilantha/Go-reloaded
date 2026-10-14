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
	anRegex     = regexp.MustCompile(`\b(?:[Aa]|[Aa]n)\b`)
)

func process(content string) string {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		line = processMarkers(line)
		line = fixQuote(line)
		line = fixPunct(line)
		line = fixArticles(line)
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

func fixArticles(line string) string {
	matches := anRegex.FindAllStringIndex(line, -1)
	if len(matches) == 0 {
		return line
	}

	var out strings.Builder
	last := 0

	for _, match := range matches {
		start, end := match[0], match[1]
		out.WriteString(line[last:start])

		article := line[start:end]
		nextWord := nextWordAfter(line, end)
		if nextWord == "" {
			out.WriteString(article)
		} else {
			out.WriteString(formatArticle(article, nextWord))
		}

		last = end
	}

	out.WriteString(line[last:])
	return out.String()
}

func nextWordAfter(line string, pos int) string {
	for pos < len(line) && unicode.IsSpace(rune(line[pos])) {
		pos++
	}
	if pos >= len(line) {
		return ""
	}

	start := pos
	for pos < len(line) && !unicode.IsSpace(rune(line[pos])) {
		pos++
	}

	word := strings.Trim(line[start:pos], `"'“”‘’.,!?:;()[]{}<>`)
	if word == "" {
		return ""
	}

	return word
}

func formatArticle(article, nextWord string) string {
	nextWord = strings.TrimFunc(nextWord, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	if nextWord == "" {
		return article
	}

	runes := []rune(strings.ToLower(nextWord))
	desired := "a"
	if isVowel(runes[0]) || (len(runes) > 1 && runes[0] == 'h' && isVowel(runes[1])) {
		desired = "an"
	}

	if article == strings.ToUpper(article) && len(article) > 1 {
		return strings.ToUpper(desired)
	}
	if article != strings.ToLower(article) {
		return strings.ToUpper(string(desired[0])) + desired[1:]
	}

	return desired
}

func isVowel(r rune) bool {
	switch unicode.ToLower(r) {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}
