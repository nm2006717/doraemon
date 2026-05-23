package search

import (
	"strings"

	"github.com/frankfoo/doraemon/internal/store"
)

type Recommender struct {
	store *store.Store
}

func NewRecommender(s *store.Store) *Recommender {
	return &Recommender{store: s}
}

func (r *Recommender) Recommend(context string, limit int) ([]store.Entry, error) {
	if limit <= 0 {
		limit = 5
	}

	keywords := extractKeywords(context)
	if len(keywords) == 0 {
		return nil, nil
	}

	query := buildRecommendQuery(keywords)
	entries, err := r.store.Search(query)
	if err != nil {
		return nil, err
	}

	if len(entries) > limit {
		entries = entries[:limit]
	}
	return entries, nil
}

func extractKeywords(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	seen := make(map[string]bool)
	var keywords []string

	for _, w := range words {
		w = strings.Trim(w, ".,;:!?\"'()[]{}")
		if len(w) < 2 || stopWords[w] || seen[w] {
			continue
		}
		seen[w] = true
		keywords = append(keywords, w)
	}

	if len(keywords) > 10 {
		keywords = keywords[:10]
	}
	return keywords
}

func buildRecommendQuery(keywords []string) string {
	parts := make([]string, len(keywords))
	for i, kw := range keywords {
		parts[i] = kw + "*"
	}
	return strings.Join(parts, " OR ")
}

var stopWords = map[string]bool{
	"the": true, "a": true, "an": true, "is": true, "are": true,
	"was": true, "were": true, "be": true, "been": true, "being": true,
	"have": true, "has": true, "had": true, "do": true, "does": true,
	"did": true, "will": true, "would": true, "could": true, "should": true,
	"may": true, "might": true, "shall": true, "can": true,
	"it": true, "its": true, "this": true, "that": true, "these": true,
	"those": true, "i": true, "you": true, "he": true, "she": true,
	"we": true, "they": true, "me": true, "him": true, "her": true,
	"us": true, "them": true, "my": true, "your": true, "his": true,
	"our": true, "their": true, "what": true, "which": true, "who": true,
	"when": true, "where": true, "why": true, "how": true,
	"not": true, "no": true, "nor": true, "but": true, "and": true,
	"or": true, "if": true, "then": true, "else": true,
	"of": true, "to": true, "in": true, "on": true, "at": true,
	"by": true, "for": true, "with": true, "from": true, "up": true,
	"about": true, "into": true, "over": true, "after": true,
	"的": true, "是": true, "了": true, "在": true, "和": true,
	"有": true, "我": true, "他": true, "她": true, "它": true,
	"们": true, "这": true, "那": true, "就": true, "也": true,
}
