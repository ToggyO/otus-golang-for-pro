package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var resultMaxLen = 10

var reg = regexp.MustCompile("[!\"#$%&'()*+,./:;<=>?@[\\]^_\\x60{|}~]|(^|\\s*)-(\\s|$)")

func Top10(text string) []string {
	raw := reg.ReplaceAllString(text, " ")
	split := strings.Fields(raw)

	var theHighestNumberOfOccurrences int
	theNumberOfOccurrencesOfEachWordMap := make(map[string]int, len(split))
	for _, word := range split {
		word := strings.ToLower(word)

		count, ok := theNumberOfOccurrencesOfEachWordMap[word]
		if !ok {
			theNumberOfOccurrencesOfEachWordMap[word] = 1
			continue
		}

		count++
		theNumberOfOccurrencesOfEachWordMap[word] = count

		if count > theHighestNumberOfOccurrences {
			theHighestNumberOfOccurrences = count
		}
	}

	// Attempt to avoid slices extra memory allocations
	// by precalculating capacities for `wordsGroupedByOccurrencesNumberMap` values
	optimizedCapacityMap := make(map[int]int, theHighestNumberOfOccurrences)
	for _, count := range theNumberOfOccurrencesOfEachWordMap {
		capacity, ok := optimizedCapacityMap[count]
		if !ok {
			optimizedCapacityMap[count] = 1
			continue
		}

		capacity++
		optimizedCapacityMap[count] = capacity
	}

	wordsGroupedByOccurrencesNumberMap := make(map[int][]string, theHighestNumberOfOccurrences)
	for k, count := range theNumberOfOccurrencesOfEachWordMap {
		words, ok := wordsGroupedByOccurrencesNumberMap[count]
		if !ok {
			c := optimizedCapacityMap[count]
			wordsGroupedByOccurrencesNumberMap[count] = append(make([]string, 0, c), k)
			continue
		}

		wordsGroupedByOccurrencesNumberMap[count] = append(words, k)
	}

	sortedByOccurrencesCount := make([]int, 0, len(wordsGroupedByOccurrencesNumberMap))
	for k := range wordsGroupedByOccurrencesNumberMap {
		sortedByOccurrencesCount = append(sortedByOccurrencesCount, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedByOccurrencesCount)))

	result := make([]string, 0, resultMaxLen)

	for _, num := range sortedByOccurrencesCount {
		words := wordsGroupedByOccurrencesNumberMap[num]
		wordsLen := len(words)
		resultLen := len(result)

		if wordsLen > 1 {
			sort.Strings(words)
		}

		if wordsLen+resultLen <= resultMaxLen {
			result = append(result, words...)
			continue
		}

		diff := resultMaxLen - resultLen
		result = append(result, words[:diff]...)
	}

	return result
}
