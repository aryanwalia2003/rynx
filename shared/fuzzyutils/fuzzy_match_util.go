package fuzzyutils

import "strings"

type Fuzzy_match_struct struct {
	Item  string
	Score int
}

func Fuzzy_score_util(item, query string) (int, bool) {
	if query == "" {
		return 0, true
	}
	item = strings.ToLower(item)
	query = strings.ToLower(query)

	score := 0
	itemIdx := 0
	queryIdx := 0
	lastMatch := -1

	for itemIdx < len(item) && queryIdx < len(query) {
		if item[itemIdx] == query[queryIdx] {
			if lastMatch == itemIdx-1 {
				score += 5 // bonus for consecutive matches
			}
			score += 1
			lastMatch = itemIdx
			queryIdx++
		}
		itemIdx++
	}

	if queryIdx != len(query) {
		return 0, false
	}
	return score, true
}

func Fuzzy_filter_util(items []string, query string) []Fuzzy_match_struct {
	results := []Fuzzy_match_struct{}
	for _, item := range items {
		if score, ok := Fuzzy_score_util(item, query); ok {
			results = append(results, Fuzzy_match_struct{Item: item, Score: score})
		}
	}
	return results
}
