package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	names, _ := readNamesFromFile()
	printFirstLetterFrequency(names)
}

func readNamesFromFile() ([]string, error) {
	var names []string
	file, err := os.Open("first-names.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return names, nil
}

func printFirstLetterFrequency(words []string) {
	count := make(map[rune]int)
	for _, word := range words {
		firstLetter := []rune(word)[0]
		count[firstLetter]++
	}
	for letter, frequency := range count {
		fmt.Printf("%c: %.0f%%\n", letter, float64(frequency)/float64(len(words))*100)
	}
}
