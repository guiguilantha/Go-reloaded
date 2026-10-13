package main

import (
	"fmt"
	"os"
)

func main() {
	inputPath, outputPath := parseArgs()
	content := readFile(inputPath)
	result := process(content)
	writeFile(outputPath, result)
}

func parseArgs() (string, string) {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: go run . <input.txt> <output.txt>")
		os.Exit(1)
	}

	return os.Args[1], os.Args[2]
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erreur de lecture:", err)
		os.Exit(1)
	}

	return string(data)
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "erreur d'écriture:", err)
		os.Exit(1)
	}
}
