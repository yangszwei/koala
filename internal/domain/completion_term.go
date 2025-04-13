package domain

// CompletionTerm represents a term for autocomplete suggestions.
type CompletionTerm struct {
	Term string `json:"term"` // The display text shown in suggestions
}

// NewCompletionTerm creates a new CompletionTerm instance.
func NewCompletionTerm(term string) CompletionTerm {
	return CompletionTerm{Term: term}
}
