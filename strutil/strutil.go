// Package strutil implements string helpers
package strutil

import "strings"

var mask = '*'

// MaskRight replaces string symbols after n with "*"
func MaskRight(s string, n int) string {
	rs := []rune(s)
	if n >= len(rs) {
		return s
	}
	for i := n; i < len(rs); i++ {
		rs[i] = mask
	}
	return string(rs)
}

// MaskLeft replaces string symbols before n with "*"
func MaskLeft(s string, n int) string {
	if n == 0 {
		return s
	}
	rs := []rune(s)
	for i := 0; i < len(rs)-n; i++ {
		rs[i] = mask
	}
	return string(rs)
}

// Decapitalize converted the first character to lower case
func Decapitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

// IsUUID checks whether given string is in uuid format
func IsUUID(s string) bool {
	if len(s) == 36 && s[8] == '-' && s[13] == '-' && s[18] == '-' && s[23] == '-' {
		return true
	}
	return false
}

// Unique returns new slice with unique elements
func Unique(keys []string) []string {
	if len(keys) == 0 {
		return keys
	}
	seen := make(map[string]struct{})
	unique := make([]string, 0)
	for _, k := range keys {
		if _, has := seen[k]; has {
			continue
		}
		seen[k] = struct{}{}
		unique = append(unique, k)
	}
	return unique
}
