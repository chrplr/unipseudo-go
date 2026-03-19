// Port of pseudoword-generation-by-markov-on-trigrams.R to Go.
// Original R script by Christophe Pallier.
//
// Copyright (c) 2026 Christophe Pallier
// This code is distributed under the terms of the GNU General Public License v3.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
	"unicode/utf8"
)

const version = "1.0.0"

type TrigramMap map[string][]rune

func generatePseudowords(n int, length int, models []string, exclude []string) []string {
	// Pad models with spaces
	paddedModels := make([][]rune, len(models))
	for i, m := range models {
		paddedModels[i] = append([]rune(" "), append([]rune(m), ' ')...)
	}

	// trigs[pos][bigram] = list of 3rd characters
	trigs := make([]map[string][]rune, length-1)
	for i := range trigs {
		trigs[i] = make(map[string][]rune)
	}

	for cpos := 0; cpos < length-1; cpos++ {
		for _, runes := range paddedModels {
			if len(runes) >= cpos+3 {
				bigram := string(runes[cpos : cpos+2])
				char3 := runes[cpos+2]
				trigs[cpos][bigram] = append(trigs[cpos][bigram], char3)
			}
		}
	}

	// Initial trigrams (from position 0)
	var initialTrigrams [][]rune
	for bigram, chars := range trigs[0] {
		for _, c := range chars {
			initialTrigrams = append(initialTrigrams, append([]rune(bigram), c))
		}
	}

	pseudos := make([]string, 0, n)
	excludeMap := make(map[string]bool)
	for _, e := range exclude {
		excludeMap[e] = true
	}
	// The R script excludes the models (padded)
	for _, runes := range paddedModels {
		excludeMap[string(runes)] = true
	}

	seen := make(map[string]bool)

	rand.Seed(time.Now().UnixNano())

	for len(pseudos) < n {
		if len(initialTrigrams) == 0 {
			break
		}
		
		itemRunes := make([]rune, 3)
		copy(itemRunes, initialTrigrams[rand.Intn(len(initialTrigrams))])

		possible := true
		for pos := 1; pos < length-1; pos++ {
			lastBigram := string(itemRunes[len(itemRunes)-2:])
			
			compatChars, ok := trigs[pos][lastBigram]
			if !ok || len(compatChars) == 0 {
				possible = false
				break
			}

			chosenChar := compatChars[rand.Intn(len(compatChars))]
			itemRunes = append(itemRunes, chosenChar)
		}

		if !possible {
			continue
		}

		finalItem := string(itemRunes)
		if !excludeMap[finalItem] && !seen[finalItem] {
			pseudos = append(pseudos, finalItem)
			seen[finalItem] = true
		}
	}

	// Remove the leading space (as per R's substring(pseudos, 2))
	results := make([]string, len(pseudos))
	for i, p := range pseudos {
		runes := []rune(p)
		if len(runes) > 0 {
			results[i] = string(runes[1:])
		}
	}

	return results
}

func main() {
	num := flag.Int("n", 10, "number of pseudowords to generate")
	length := flag.Int("l", 7, "length of the pseudowords")
	minLen := flag.Int("m", 5, "minimum length of model words")
	inputFile := flag.String("f", "liste.de.mots.francais.frgut.txt", "input word list file")
	showVersion := flag.Bool("v", false, "display version, license, author, and source location")
	flag.Parse()

	if *showVersion {
		fmt.Printf("unipseudo-go version %s\n", version)
		fmt.Println("Author: Christophe Pallier (2026)")
		fmt.Println("License: GNU General Public License v3")
		fmt.Println("Source: http://github.com/chrplr/unipseudo-go")
		os.Exit(0)
	}

	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var models []string
	for _, w := range words {
		if utf8.RuneCountInString(w) > *minLen {
			models = append(models, w)
		}
	}

	pseudos := generatePseudowords(*num, *length, models, nil)
	for _, p := range pseudos {
		fmt.Println(p)
	}
}
