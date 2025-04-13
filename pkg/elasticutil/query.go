package elasticutil

import "regexp"

// EscapeQueryString escapes special characters for Elasticsearch queries.
func EscapeQueryString(s string) string {
	re := regexp.MustCompile(`([+\-=&|!(){}\[\]^"~*?:\\/<>])`)
	return re.ReplaceAllString(s, `\\$1`)
}
