package hw09structvalidator

import "strings"

func parseTag(tag string) []validationCondition {
	splitByAnd := strings.Split(tag, separatorAnd)
	resultLen := len(splitByAnd)
	result := make([]validationCondition, resultLen)

	for i := 0; i < resultLen; i++ {
		c := strings.Split(splitByAnd[i], separatorTagValue)
		result[i] = validationCondition{name: c[0], rule: c[1]}
	}

	return result
}
