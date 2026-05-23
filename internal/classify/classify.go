package classify

import (
	"net/url"
	"path/filepath"
	"strings"
)

type Result struct {
	Category string
	Tags     []string
}

func Classify(title, content string) Result {
	category := detectCategory(content)
	tags := extractTags(title, content)
	return Result{Category: category, Tags: tags}
}

func detectCategory(content string) string {
	trimmed := strings.TrimSpace(content)

	if looksLikeURL(trimmed) {
		return "link"
	}

	if looksLikeImage(trimmed) {
		return "image"
	}

	if looksLikeFile(trimmed) {
		return "file"
	}

	return "note"
}

func looksLikeURL(content string) bool {
	lines := strings.Split(content, "\n")
	first := strings.TrimSpace(lines[0])
	u, err := url.Parse(first)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func looksLikeImage(content string) bool {
	lower := strings.ToLower(content)
	imageExts := []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".bmp"}
	for _, ext := range imageExts {
		if strings.Contains(lower, ext) && (strings.Contains(lower, "image") || strings.HasSuffix(strings.TrimSpace(lower), ext)) {
			return true
		}
	}
	return false
}

func looksLikeFile(content string) bool {
	lower := strings.ToLower(strings.TrimSpace(content))
	fileExts := []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".csv", ".zip", ".tar"}
	for _, ext := range fileExts {
		if strings.HasSuffix(lower, ext) || strings.Contains(lower, ext) {
			if filepath.Ext(lower) == ext {
				return true
			}
		}
	}
	return false
}

var tagKeywords = map[string][]string{
	"go":         {"go", "golang", "goroutine", "channel", "defer"},
	"python":     {"python", "pip", "django", "flask", "pandas"},
	"javascript": {"javascript", "js", "node", "react", "vue", "typescript", "ts"},
	"rust":       {"rust", "cargo", "crate"},
	"docker":     {"docker", "container", "dockerfile", "compose"},
	"k8s":        {"kubernetes", "k8s", "kubectl", "pod", "deployment"},
	"database":   {"sql", "sqlite", "postgres", "mysql", "mongodb", "redis", "database", "db"},
	"api":        {"api", "rest", "grpc", "graphql", "endpoint"},
	"git":        {"git", "github", "gitlab", "branch", "merge", "commit"},
	"testing":    {"test", "testing", "unittest", "benchmark"},
	"concurrency": {"goroutine", "channel", "mutex", "async", "await", "thread", "concurrent"},
	"web":        {"http", "html", "css", "web", "browser", "frontend"},
	"devops":     {"ci", "cd", "pipeline", "deploy", "terraform", "ansible"},
	"security":   {"auth", "token", "jwt", "oauth", "encryption", "security"},
	"ai":         {"ai", "ml", "llm", "gpt", "model", "neural", "machine learning"},
}

func extractTags(title, content string) []string {
	text := strings.ToLower(title + " " + content)
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'))
	})

	wordSet := make(map[string]bool, len(words))
	for _, w := range words {
		wordSet[w] = true
	}

	matched := make(map[string]bool)
	for tag, keywords := range tagKeywords {
		for _, kw := range keywords {
			if wordSet[kw] {
				matched[tag] = true
				break
			}
		}
	}

	tags := make([]string, 0, len(matched))
	for tag := range matched {
		tags = append(tags, tag)
	}
	return tags
}

func FormatTags(tags []string) string {
	return strings.Join(tags, ",")
}
