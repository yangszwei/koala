package elasticutil

import "regexp"

// EscapeQueryString escapes special characters for Elasticsearch queries.
func EscapeQueryString(s string) string {
	re := regexp.MustCompile(`([+\-=&|!(){}\[\]^"~*?:\\/<>])`)
	return re.ReplaceAllString(s, `\\$1`)
}

// EscapeQueryStrings escapes special characters in a slice of strings.
func EscapeQueryStrings(ss []string) []string {
	escaped := make([]string, len(ss))
	for i, s := range ss {
		escaped[i] = EscapeQueryString(s)
	}
	return escaped
}
