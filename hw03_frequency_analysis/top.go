package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"strings"
	"unicode"
)

// MapWordFrequency split str and return map, where key is word and value is frequency.
func MapWordFrequency(str string) map[string]int {
	wordFrequency := make(map[string]int)

	words := strings.FieldsFunc(str, func(r rune) bool {
		// rune ' for English words(like It's) with an apostrophe
		return !unicode.IsLetter(r) && r != '-' && r != '\''
	})

	for _, w := range words {
		if w != "-" {
			w = strings.ToLower(w)
			w = strings.Trim(w, "'")
			wordFrequency[w]++
		}
	}

	return wordFrequency
}

// MapFrequencyWords convert words map where key is word and value is frequecy
// to map where key is frequency and value is slice of words and return it
// also return max frequncy.
func MapFrequencyWords(words map[string]int) (map[int][]string, int) {
	frequencyWords, maxFrequency := make(map[int][]string), 0

	for word, freq := range words {
		if _, ok := frequencyWords[freq]; !ok {
			frequencyWords[freq] = make([]string, 0)
		}
		frequencyWords[freq] = append(frequencyWords[freq], word)

		if freq > maxFrequency {
			maxFrequency = freq
		}
	}

	return frequencyWords, maxFrequency
}

func Top10(str string) []string {
	words := MapWordFrequency(str)
	frequencyWords, maxFrequency := MapFrequencyWords(words)

	const nMostFrequency = 10
	mostFrequency := make([]string, 0, nMostFrequency)

	for i := maxFrequency; i > 0; i-- {
		if _, ok := frequencyWords[i]; !ok {
			continue
		}

		if len(mostFrequency)+len(frequencyWords[i]) > nMostFrequency {
			last := nMostFrequency - len(mostFrequency)
			mostFrequency = append(mostFrequency, frequencyWords[i][:last]...)
			break
		}

		mostFrequency = append(mostFrequency, frequencyWords[i]...)
	}

	return mostFrequency
}
