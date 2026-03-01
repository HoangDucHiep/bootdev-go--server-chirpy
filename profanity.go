package main

import (
	"strings"
	"unicode"
)

func CleanProfanity(body string) string {
	profaneWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	words := strings.Split(body, " ")
	for i, word := range words {
		cleanedWord, prefix, suffix := extractWord(word)
		if profaneWords[strings.ToLower(cleanedWord)] {
			words[i] = prefix + "****" + suffix
		}
	}

	return strings.Join(words, " ")
}

func extractWord(word string) (cleanedWord, prefix, suffix string) {
	runes := []rune(word)

	start := 0
	for start < len(runes) && !unicode.IsLetter(runes[start]) {
		start++
	}

	end := len(runes)
	for end > start && !unicode.IsLetter(runes[end-1]) {
		end--
	}

	prefix = string(runes[:start])
	cleanedWord = string(runes[start:end])
	suffix = string(runes[end:])

	return cleanedWord, prefix, suffix
}
