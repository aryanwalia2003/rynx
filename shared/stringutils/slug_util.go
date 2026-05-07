package stringutils

import (
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// Slugify_util converts a string into a URL/git friendly slug
func Slugify_util(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAlphanumericRegex.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	
	// Limit length to something reasonable to avoid OS/git issues
	if len(s) > 50 {
		s = s[:50]
		s = strings.TrimRight(s, "-")
	}
	
	return s
}
