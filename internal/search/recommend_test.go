package search

import (
	"testing"
)

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			"filters stop words",
			"the quick brown fox is running",
			[]string{"quick", "brown", "fox", "running"},
		},
		{
			"filters chinese stop words",
			"我 在 学习 golang 并发",
			[]string{"学习", "golang", "并发"},
		},
		{
			"removes punctuation",
			"hello, world! (test)",
			[]string{"hello", "world", "test"},
		},
		{
			"deduplicates",
			"go go go golang golang",
			[]string{"go", "golang"},
		},
		{
			"skips short words",
			"a b c de fg",
			[]string{"de", "fg"},
		},
		{
			"caps at 10",
			"one two three four five six seven eight nine ten eleven twelve",
			[]string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"},
		},
		{
			"empty input",
			"",
			nil,
		},
		{
			"only stop words",
			"the is are was were",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractKeywords(tt.text)
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("extractKeywords(%q) = %v, want %v", tt.text, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("extractKeywords(%q)[%d] = %q, want %q", tt.text, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestBuildRecommendQuery(t *testing.T) {
	tests := []struct {
		name     string
		keywords []string
		want     string
	}{
		{"single", []string{"golang"}, "golang*"},
		{"multiple", []string{"go", "channel"}, "go* OR channel*"},
		{"three", []string{"docker", "k8s", "deploy"}, "docker* OR k8s* OR deploy*"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRecommendQuery(tt.keywords)
			if got != tt.want {
				t.Errorf("buildRecommendQuery(%v) = %q, want %q", tt.keywords, got, tt.want)
			}
		})
	}
}
